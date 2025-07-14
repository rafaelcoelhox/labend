package app

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/rafaelcoelhox/labbend/internal/challenges"
	schemas_configuration "github.com/rafaelcoelhox/labbend/internal/config/graphql"
	"github.com/rafaelcoelhox/labbend/internal/users"
	"github.com/rafaelcoelhox/labbend/pkg/database"
	"github.com/rafaelcoelhox/labbend/pkg/eventbus"
	"github.com/rafaelcoelhox/labbend/pkg/health"
	corelogger "github.com/rafaelcoelhox/labbend/pkg/logger"
	"github.com/rafaelcoelhox/labbend/pkg/monitoring"
	"github.com/rafaelcoelhox/labbend/pkg/saga"
	"gorm.io/gorm/logger"
)

// App - estrutura principal da aplicação
type App struct {
	config      Config
	db          *gorm.DB
	logger      corelogger.Logger
	eventBus    *eventbus.EventBus
	txManager   *database.TxManager
	sagaManager *saga.SagaManager
	healthMgr   *health.Manager
	monitor     *monitoring.Monitor
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
	}

	log, err = corelogger.NewWithConfig(loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Configurar logger do GORM
	var gormLogLevel logger.LogLevel
	switch config.LogLevel {
	case "debug":
		gormLogLevel = logger.Info
	case "warn":
		gormLogLevel = logger.Warn
	case "error":
		gormLogLevel = logger.Error
	default:
		gormLogLevel = logger.Silent
	}

	// Setup database
	db, err := database.Connect(database.Config{
		DSN:          config.DatabaseURL,
		MaxIdleConns: config.MaxIdleConns,
		MaxOpenConns: config.MaxOpenConns,
		MaxLifetime:  config.ConnMaxLifetime,
		LogLevel:     gormLogLevel,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Info("Database connection established", zap.String("database", config.DatabaseURL))

	// Auto migrate database tables using registered models
	registeredModels := database.GetRegisteredModels()
	log.Info("Auto migrating database", zap.Int("registered_models", len(registeredModels)))

	if err := database.AutoMigrateRegistered(db); err != nil {
		return nil, fmt.Errorf("failed to auto migrate database: %w", err)
	}
	log.Info("Database auto migration completed")

	// Setup database transaction manager
	txManager := database.NewTxManager(db)

	// Setup event bus
	eventBus := eventbus.New(log)

	// Setup saga manager
	sagaManager := saga.NewSagaManager(log)

	// Setup health manager
	healthMgr := health.NewManager()
	healthMgr.Register("database", health.NewDatabaseChecker(db))

	// Setup monitoring
	monitor := monitoring.NewMonitor(log)

	return &App{
		config:      config,
		db:          db,
		logger:      log,
		eventBus:    eventBus,
		txManager:   txManager,
		sagaManager: sagaManager,
		healthMgr:   healthMgr,
		monitor:     monitor,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	a.logger.Info("Starting application", zap.String("environment", a.config.Environment))

	// Setup repositories
	userRepo := users.NewRepository(a.db)
	challengeRepo := challenges.NewRepository(a.db)

	// Setup services
	userService := users.NewService(userRepo, a.logger, a.eventBus, a.txManager)
	challengeService := challenges.NewService(challengeRepo, userService, a.logger, a.eventBus, a.txManager, a.sagaManager)

	// Setup GraphQL schema usando o novo ModuleRegistry
	registry := schemas_configuration.NewModuleRegistry(a.logger)
	registry.Register("users", userService)
	registry.Register("challenges", challengeService)
	// Adicione novos módulos aqui: registry.Register("products", productService)

	schema, err := schemas_configuration.ConfigureSchema(registry)
	if err != nil {
		return fmt.Errorf("failed to build GraphQL schema: %w", err)
	}

	// Setup server
	if a.config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware básico
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		status := a.healthMgr.CheckAll(context.Background())
		statusCode := http.StatusOK
		if status.Status != health.StatusHealthy {
			statusCode = http.StatusServiceUnavailable
		}
		c.JSON(statusCode, status)
	})

	// Metrics endpoint
	router.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "metrics endpoint"})
	})

	// Middleware de CORS simples
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Setup GraphQL handler usando graphql-go/handler
	graphqlHandler := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     !a.config.IsProduction(), // JSON formatado apenas em desenvolvimento
		GraphiQL:   !a.config.IsProduction(), // Interface GraphiQL apenas em desenvolvimento
		Playground: false,
	})

	// GraphQL endpoint
	router.POST("/graphql", func(c *gin.Context) {
		graphqlHandler.ServeHTTP(c.Writer, c.Request)
	})

	// GraphQL playground (apenas em desenvolvimento)
	if !a.config.IsProduction() {
		router.GET("/graphql", func(c *gin.Context) {
			graphqlHandler.ServeHTTP(c.Writer, c.Request)
		})
	}

	// Health check simples
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "LabEnd API",
			"version":     "1.0.0",
			"environment": a.config.Environment,
			"status":      "healthy",
		})
	})

	// Configurar timeouts do servidor
	port, _ := strconv.Atoi(a.config.Port)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	a.logger.Info("Server starting", zap.String("port", a.config.Port))

	// Iniciar servidor em goroutine separada
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Aguardar sinal de shutdown
	<-ctx.Done()

	a.logger.Info("Shutting down server...")

	// Context com timeout para o shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("Server forced to shutdown", zap.Error(err))
		return err
	}

	a.logger.Info("Server exited properly")
	return nil
}

func (a *App) Stop() error {
	a.logger.Info("Application stopping...")

	// Fechar event bus
	if a.eventBus != nil {
		a.eventBus.Shutdown()
	}
	// Fechar conexão com o banco de dados
	if sqlDB, err := a.db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			a.logger.Error("Failed to close database connection", zap.Error(err))
			return err
		}
	}

	a.logger.Info("Application stopped successfully")
	return nil
}

// userServiceAdapter adapta o users.Service para ser compatível com outros módulos
type userServiceAdapter struct {
	userService users.Service
}

func (u *userServiceAdapter) GetUser(ctx context.Context, id uint) (interface{}, error) {
	return u.userService.GetUser(ctx, id)
}
