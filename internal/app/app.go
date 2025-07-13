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
	"github.com/rafaelcoelhox/labbend/internal/core/logger"
	"github.com/rafaelcoelhox/labbend/internal/users"
)

// App - aplicação principal
type App struct {
	config    Config
	db        *gorm.DB
	logger    logger.Logger
	eventBus  *eventbus.EventBus
	healthMgr *health.Manager
	server    *http.Server
}

// New - cria nova instância da aplicação
func New(config Config) (*App, error) {
	// Setup logger
	log, err := logger.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Setup database
	dbConfig := database.DefaultConfig(config.DatabaseURL)
	db, err := database.Connect(dbConfig)
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

	// Setup event bus
	eventBus := eventbus.New(log)

	// Setup health manager
	healthMgr := health.NewManager()
	healthMgr.Register("database", health.NewDatabaseChecker(db))
	healthMgr.Register("memory", health.NewMemoryChecker())
	healthMgr.Register("eventbus", health.NewEventBusChecker(eventBus))

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

	// Setup HTTP server
	router := gin.Default()

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
			// Usar método otimizado para buscar usuários com XP
			users, err := userResolver.Users(c.Request.Context(), nil, nil)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, users)
		})

		api.GET("/users/:id", func(c *gin.Context) {
			id := c.Param("id")
			user, err := userResolver.User(c.Request.Context(), id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, user)
		})

		api.POST("/users", func(c *gin.Context) {
			var input users.CreateUserInput
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			user, err := userResolver.CreateUser(c.Request.Context(), input)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, user)
		})

		// Challenge routes
		api.GET("/challenges", func(c *gin.Context) {
			challenges, err := challengeResolver.Challenges(c.Request.Context(), nil, nil)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, challenges)
		})

		api.POST("/challenges", func(c *gin.Context) {
			var input challenges.CreateChallengeInput
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			challenge, err := challengeResolver.CreateChallenge(c.Request.Context(), input)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, challenge)
		})
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Detailed health check
	router.GET("/health/detailed", func(c *gin.Context) {
		report := a.healthMgr.CheckAll(c.Request.Context())

		status := http.StatusOK
		if report.Status == health.StatusUnhealthy {
			status = http.StatusServiceUnavailable
		} else if report.Status == health.StatusDegraded {
			status = http.StatusPartialContent
		}

		c.JSON(status, report)
	})

	// Metrics endpoint
	router.GET("/metrics", func(c *gin.Context) {
		// Em produção usaria Prometheus metrics
		c.JSON(http.StatusOK, gin.H{
			"uptime_seconds": time.Since(time.Now().Add(-1 * time.Hour)).Seconds(), // Mock
			"modules": gin.H{
				"users":      gin.H{"status": "active"},
				"challenges": gin.H{"status": "active"},
			},
		})
	})

	a.server = &http.Server{
		Addr:           ":" + a.config.Port,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	a.logger.Info("starting server", zap.String("port", a.config.Port))

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

	a.logger.Info("shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown event bus
	a.eventBus.Shutdown()

	// Shutdown HTTP server
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("server forced to shutdown", zap.Error(err))
		return err
	}

	a.logger.Info("server exited")
	return nil
}
