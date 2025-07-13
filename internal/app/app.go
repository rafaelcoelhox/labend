package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/rafaelcoelhox/labbend/internal/app/graph"
	"github.com/rafaelcoelhox/labbend/internal/challenges"
	"github.com/rafaelcoelhox/labbend/internal/core/database"
	"github.com/rafaelcoelhox/labbend/internal/core/eventbus"
	"github.com/rafaelcoelhox/labbend/internal/core/health"
	corelogger "github.com/rafaelcoelhox/labbend/internal/core/logger"
	"github.com/rafaelcoelhox/labbend/internal/core/saga"
	"github.com/rafaelcoelhox/labbend/internal/users"
	"gorm.io/gorm/logger"
)

// App - estrutura principal da aplicação
type App struct {
	config           Config
	db               *gorm.DB
	logger           corelogger.Logger
	eventBusManager  *eventbus.EventBusManager
	txManager        *database.TxManager
	sagaManager      *saga.SagaManager
	healthMgr        *health.Manager
	outboxRepository *eventbus.OutboxRepository
}

// NewApp - cria nova instância da aplicação
func NewApp(config Config) (*App, error) {
	// Setup logger baseado no ambiente
	var log corelogger.Logger
	var err error

	loggerConfig := corelogger.Config{
		Level:            config.LogLevel,
		Environment:      config.Environment,
		EnableCaller:     true,
		EnableStacktrace: config.IsProduction(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	log, err = corelogger.NewWithConfig(loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Setup database com configuração avançada
	dbConfig := config.GetDatabaseConfig()

	// Configurar log level baseado no ambiente
	var logLevel logger.LogLevel
	if config.IsProduction() {
		logLevel = logger.Warn
	} else {
		logLevel = logger.Info
	}

	db, err := database.Connect(database.Config{
		DSN:          dbConfig.DSN,
		MaxIdleConns: dbConfig.MaxIdleConns,
		MaxOpenConns: dbConfig.MaxOpenConns,
		MaxLifetime:  dbConfig.ConnMaxLifetime,
		LogLevel:     logLevel,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Setup Transaction Manager
	txManager := database.NewTxManager(db)

	// Setup Saga Manager
	sagaManager := saga.NewSagaManager(log)

	// Setup Event Bus components
	immediateEventBus := eventbus.New(log)
	outboxRepository := eventbus.NewOutboxRepository(db)
	transactionalEventBus := eventbus.NewTransactionalEventBus(immediateEventBus, outboxRepository, log)
	eventBusManager := eventbus.NewEventBusManager(immediateEventBus, transactionalEventBus, log)

	// Auto migrate (incluindo tabela de outbox)
	if err := database.AutoMigrate(db,
		&users.User{},
		&users.UserXP{},
		&challenges.Challenge{},
		&challenges.ChallengeSubmission{},
		&challenges.ChallengeVote{},
		&eventbus.OutboxEvent{}, // Tabela do outbox
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Setup health manager
	healthMgr := health.NewManager()
	healthMgr.Register("database", health.NewDatabaseChecker(db))
	healthMgr.Register("memory", health.NewMemoryChecker())
	healthMgr.Register("eventbus", health.NewEventBusChecker(immediateEventBus))

	log.Info("application initialized",
		zap.String("environment", config.Environment),
		zap.String("port", config.Port),
		zap.Bool("production", config.IsProduction()),
		zap.Bool("transactional_events", true),
		zap.Bool("saga_enabled", true),
	)

	return &App{
		config:           config,
		db:               db,
		logger:           log,
		eventBusManager:  eventBusManager,
		txManager:        txManager,
		sagaManager:      sagaManager,
		healthMgr:        healthMgr,
		outboxRepository: outboxRepository,
	}, nil
}

// Start - inicia a aplicação
func (a *App) Start(ctx context.Context) error {
	a.logger.Info("starting application")

	// Iniciar Event Bus Manager (processamento de outbox em background)
	a.eventBusManager.Start(ctx)

	// Setup repositories
	userRepo := users.NewRepository(a.db)
	challengeRepo := challenges.NewRepository(a.db)

	// Setup services com novos componentes
	userService := users.NewService(userRepo, a.logger, a.eventBusManager.GetTransactional(), a.txManager)
	challengeService := challenges.NewService(challengeRepo, userService, a.logger, a.eventBusManager.GetTransactional(), a.txManager, a.sagaManager)

	// Setup resolvers
	userResolver := users.NewResolver(userService, a.logger)
	challengeResolver := challenges.NewResolver(challengeService, a.logger)

	// Setup GraphQL resolver
	graphqlResolver := graph.NewResolver(userService, challengeService, a.logger)

	// Setup HTTP server
	gin.SetMode(gin.ReleaseMode)
	if !a.config.IsProduction() {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	// Middleware de logging
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		// Log da requisição
		a.logger.HTTP(method, path, c.Writer.Status(), time.Since(start),
			corelogger.String("client_ip", c.ClientIP()),
			corelogger.String("user_agent", c.Request.UserAgent()),
			corelogger.Int("body_size", c.Writer.Size()),
		)
	})

	// Recovery middleware com logging customizado
	router.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				a.logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				c.Abort()
			}
		}()
		c.Next()
	})

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// GraphQL endpoint
	schema := graph.NewExecutableSchema(graph.Config{Resolvers: graphqlResolver})
	graphqlHandler := handler.NewDefaultServer(schema)
	playgroundHandler := playground.Handler("GraphQL Playground", "/graphql")

	router.POST("/graphql", gin.WrapH(graphqlHandler))
	router.GET("/graphql", gin.WrapH(playgroundHandler))

	// API routes
	api := router.Group("/api")
	{
		// User routes
		api.GET("/users", func(c *gin.Context) {
			start := time.Now()
			// Usar método otimizado para buscar usuários com XP
			users, err := userResolver.Users(c.Request.Context(), nil, nil)
			if err != nil {
				a.logger.Error("failed to get users",
					zap.Error(err),
					zap.String("endpoint", "GET /api/users"),
				)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			a.logger.Performance("get_users", time.Since(start),
				zap.Int("user_count", len(users)),
			)
			c.JSON(http.StatusOK, users)
		})

		api.GET("/users/:id", func(c *gin.Context) {
			id := c.Param("id")
			requestLogger := a.logger.WithFields(zap.String("user_id", id))

			user, err := userResolver.User(c.Request.Context(), id)
			if err != nil {
				requestLogger.Error("failed to get user", zap.Error(err))
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			requestLogger.Debug("user retrieved successfully")
			c.JSON(http.StatusOK, user)
		})

		api.POST("/users", func(c *gin.Context) {
			var input users.CreateUserInput
			if err := c.ShouldBindJSON(&input); err != nil {
				a.logger.Warn("invalid user input",
					zap.Error(err),
					zap.Any("input", input),
				)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			start := time.Now()
			user, err := userResolver.CreateUser(c.Request.Context(), input)
			if err != nil {
				a.logger.Error("failed to create user",
					zap.Error(err),
					zap.String("email", input.Email),
				)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			a.logger.Event("user_created", "api",
				zap.String("user_email", user.Email),
				zap.Uint("user_id", user.ID),
			)
			a.logger.Performance("create_user", time.Since(start))
			c.JSON(http.StatusCreated, user)
		})

		// Challenge routes
		api.GET("/challenges", func(c *gin.Context) {
			start := time.Now()
			challenges, err := challengeResolver.Challenges(c.Request.Context(), nil, nil)
			if err != nil {
				a.logger.Error("failed to get challenges",
					zap.Error(err),
					zap.String("endpoint", "GET /api/challenges"),
				)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			a.logger.Performance("get_challenges", time.Since(start),
				zap.Int("challenge_count", len(challenges)),
			)
			c.JSON(http.StatusOK, challenges)
		})

		api.POST("/challenges", func(c *gin.Context) {
			var input challenges.CreateChallengeInput
			if err := c.ShouldBindJSON(&input); err != nil {
				a.logger.Warn("invalid challenge input",
					zap.Error(err),
					zap.Any("input", input),
				)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			start := time.Now()
			challenge, err := challengeResolver.CreateChallenge(c.Request.Context(), input)
			if err != nil {
				a.logger.Error("failed to create challenge",
					zap.Error(err),
					zap.String("title", input.Title),
				)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			a.logger.Event("challenge_created", "api",
				zap.String("challenge_title", challenge.Title),
				zap.Uint("challenge_id", challenge.ID),
				zap.Int("xp_reward", challenge.XPReward),
			)
			a.logger.Performance("create_challenge", time.Since(start))
			c.JSON(http.StatusCreated, challenge)
		})
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		a.logger.Debug("health check requested")
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Health check detalhado
	router.GET("/health/detailed", func(c *gin.Context) {
		a.logger.Debug("detailed health check requested")

		// Verificar status do outbox
		outboxStats, err := a.eventBusManager.GetTransactional().GetOutboxStats(ctx)
		if err != nil {
			a.logger.Error("failed to get outbox stats", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get outbox stats"})
			return
		}

		// Verificar sagas em execução
		runningSagas := a.sagaManager.GetRunningSagas()

		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"timestamp":   time.Now().Format(time.RFC3339),
			"environment": a.config.Environment,
			"outbox": gin.H{
				"pending_events": outboxStats.PendingEvents,
				"failed_events":  outboxStats.FailedEvents,
			},
			"sagas": gin.H{
				"running_count": len(runningSagas),
			},
			"database": gin.H{
				"connected": true,
			},
		})
	})

	// Endpoint para estatísticas do outbox
	router.GET("/admin/outbox/stats", func(c *gin.Context) {
		stats, err := a.eventBusManager.GetTransactional().GetOutboxStats(ctx)
		if err != nil {
			a.logger.Error("failed to get outbox stats", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, stats)
	})

	// Endpoint para forçar processamento do outbox
	router.POST("/admin/outbox/process", func(c *gin.Context) {
		if err := a.eventBusManager.GetTransactional().ProcessOutboxEvents(ctx); err != nil {
			a.logger.Error("failed to process outbox events", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "outbox processing triggered"})
	})

	// Endpoint para sagas em execução
	router.GET("/admin/sagas", func(c *gin.Context) {
		runningSagas := a.sagaManager.GetRunningSagas()

		response := make(map[string]interface{})
		for id, saga := range runningSagas {
			response[id] = gin.H{
				"progress":       saga.GetProgress(),
				"executed_steps": saga.GetExecutedSteps(),
				"total_steps":    saga.GetTotalSteps(),
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"running_sagas": response,
			"total_count":   len(runningSagas),
		})
	})

	// Metrics endpoint
	router.GET("/metrics", func(c *gin.Context) {
		a.logger.Debug("metrics requested")
		// Em produção usaria Prometheus metrics
		c.JSON(http.StatusOK, gin.H{
			"uptime_seconds": time.Since(time.Now().Add(-1 * time.Hour)).Seconds(), // Mock
			"environment":    a.config.Environment,
			"log_level":      a.config.LogLevel,
			"modules": gin.H{
				"users":      gin.H{"status": "active"},
				"challenges": gin.H{"status": "active"},
			},
		})
	})

	// Configurar server HTTP
	server := &http.Server{
		Addr:           ":" + a.config.Port,
		Handler:        router,
		ReadTimeout:    a.config.ReadTimeout,
		WriteTimeout:   a.config.WriteTimeout,
		IdleTimeout:    a.config.IdleTimeout,
		MaxHeaderBytes: a.config.MaxHeaderBytes,
	}

	a.logger.Info("server starting",
		zap.String("address", server.Addr),
		zap.Duration("read_timeout", a.config.ReadTimeout),
		zap.Duration("write_timeout", a.config.WriteTimeout),
	)

	// Iniciar servidor
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Stop - para a aplicação gracefully
func (a *App) Stop() error {
	a.logger.Info("stopping application")

	// Parar Event Bus Manager
	a.eventBusManager.Shutdown()

	a.logger.Info("application stopped")
	return nil
}
