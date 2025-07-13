package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/rafaelcoelhox/labbend/internal/challenges"
	"github.com/rafaelcoelhox/labbend/internal/core/database"
	"github.com/rafaelcoelhox/labbend/internal/core/eventbus"
	"github.com/rafaelcoelhox/labbend/internal/core/health"
	corelogger "github.com/rafaelcoelhox/labbend/internal/core/logger"
	"github.com/rafaelcoelhox/labbend/internal/users"
	"gorm.io/gorm/logger"
)

// App - aplicação principal
type App struct {
	config    Config
	db        *gorm.DB
	logger    corelogger.Logger
	eventBus  *eventbus.EventBus
	healthMgr *health.Manager
	server    *http.Server
}

// New - cria nova instância da aplicação
func New(config Config) (*App, error) {
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

	// Auto migrate
	if err := database.AutoMigrate(db,
		&users.User{},
		&users.UserXP{},
		&challenges.Challenge{},
		&challenges.ChallengeSubmission{},
		&challenges.ChallengeVote{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Setup event bus com configuração avançada
	eventBus := eventbus.New(log)

	// Setup health manager
	healthMgr := health.NewManager()
	healthMgr.Register("database", health.NewDatabaseChecker(db))
	healthMgr.Register("memory", health.NewMemoryChecker())
	healthMgr.Register("eventbus", health.NewEventBusChecker(eventBus))

	log.Info("application initialized",
		zap.String("environment", config.Environment),
		zap.String("port", config.Port),
		zap.Bool("production", config.IsProduction()),
	)

	return &App{
		config:    config,
		db:        db,
		logger:    log,
		eventBus:  eventBus,
		healthMgr: healthMgr,
	}, nil
}

// Start - inicia a aplicação
func (a *App) Start() error {
	// Setup modules
	userRepo := users.NewRepository(a.db)
	userService := users.NewService(userRepo, a.logger, a.eventBus)
	userResolver := users.NewResolver(userService, a.logger)

	challengeRepo := challenges.NewRepository(a.db)
	challengeService := challenges.NewService(challengeRepo, userService, a.logger, a.eventBus)
	challengeResolver := challenges.NewResolver(challengeService, a.logger)

	// Setup HTTP server com configuração do ambiente
	if a.config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New() // Usar New() ao invés de Default() para controle total

	// Middleware de logging HTTP customizado
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log após processar
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		// Usar o novo método HTTP do logger
		a.logger.HTTP(method, path, statusCode, latency,
			zap.String("client_ip", clientIP),
			zap.Int("body_size", c.Writer.Size()),
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

	// Detailed health check
	router.GET("/health/detailed", func(c *gin.Context) {
		start := time.Now()
		report := a.healthMgr.CheckAll(c.Request.Context())
		duration := time.Since(start)

		status := http.StatusOK
		if report.Status == health.StatusUnhealthy {
			status = http.StatusServiceUnavailable
			a.logger.Error("health check failed",
				zap.Any("report", report),
				zap.Duration("check_duration", duration),
			)
		} else if report.Status == health.StatusDegraded {
			status = http.StatusPartialContent
			a.logger.Warn("health check degraded",
				zap.Any("report", report),
				zap.Duration("check_duration", duration),
			)
		} else {
			a.logger.Debug("health check passed",
				zap.Duration("check_duration", duration),
			)
		}

		c.JSON(status, report)
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

	// Usar configurações avançadas do servidor
	a.server = &http.Server{
		Addr:           ":" + a.config.Port,
		Handler:        router,
		ReadTimeout:    a.config.ReadTimeout,
		WriteTimeout:   a.config.WriteTimeout,
		IdleTimeout:    a.config.IdleTimeout,
		MaxHeaderBytes: a.config.MaxHeaderBytes,
	}

	a.logger.Info("starting server",
		zap.String("port", a.config.Port),
		zap.Duration("read_timeout", a.config.ReadTimeout),
		zap.Duration("write_timeout", a.config.WriteTimeout),
		zap.Duration("idle_timeout", a.config.IdleTimeout),
	)

	// Start server in goroutine
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("shutdown signal received, starting graceful shutdown...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown event bus
	a.logger.Info("shutting down event bus...")
	a.eventBus.Shutdown()

	// Shutdown HTTP server
	a.logger.Info("shutting down HTTP server...")
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("server forced to shutdown", zap.Error(err))
		return err
	}

	// Sync logger
	if err := a.logger.Sync(); err != nil {
		// Don't return error, just log it
		fmt.Printf("Failed to sync logger: %v\n", err)
	}

	a.logger.Info("server shutdown completed gracefully")
	return nil
}
