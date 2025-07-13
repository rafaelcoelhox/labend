package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"ecommerce/internal/challenges"
	"ecommerce/internal/core/database"
	"ecommerce/internal/core/eventbus"
	"ecommerce/internal/core/health"
	"ecommerce/internal/core/logger"
	"ecommerce/internal/users"
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

	// Users module
	userRepo := users.NewRepository(a.db)
	userService := users.NewService(userRepo, a.logger, a.eventBus)
	userResolver := users.NewResolver(userService, a.logger)

	// Challenges module
	challengeRepo := challenges.NewRepository(a.db)
	challengeService := challenges.NewService(challengeRepo, userService, a.logger, a.eventBus)
	challengeResolver := challenges.NewResolver(challengeService, a.logger)

	// Setup GraphQL
	resolver := &Resolver{
		users:      userResolver,
		challenges: challengeResolver,
	}

	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))

	// Setup HTTP server
	router := gin.Default()

	// GraphQL endpoint
	router.POST("/graphql", func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	})

	// GraphQL playground
	router.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL Playground", "/graphql").ServeHTTP(c.Writer, c.Request)
	})

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
		Addr:    ":" + a.config.Port,
		Handler: router,
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("server forced to shutdown", zap.Error(err))
		return err
	}

	a.logger.Info("server exited")
	return nil
}

// Resolver - root resolver
type Resolver struct {
	users      *users.Resolver
	challenges *challenges.Resolver
}

// Query resolver
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// Mutation resolver
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// queryResolver - implementa QueryResolver
type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context, id string) (*users.GraphQLUser, error) {
	return r.users.User(ctx, id)
}

func (r *queryResolver) Users(ctx context.Context, limit *int, offset *int) ([]*users.GraphQLUser, error) {
	return r.users.Users(ctx, limit, offset)
}

func (r *queryResolver) UserXPHistory(ctx context.Context, userID string) ([]*users.UserXP, error) {
	return r.users.UserXPHistory(ctx, userID)
}

func (r *queryResolver) Challenge(ctx context.Context, id string) (*challenges.Challenge, error) {
	return r.challenges.Challenge(ctx, id)
}

func (r *queryResolver) Challenges(ctx context.Context, limit *int, offset *int) ([]*challenges.Challenge, error) {
	return r.challenges.Challenges(ctx, limit, offset)
}

func (r *queryResolver) ChallengeSubmissions(ctx context.Context, challengeID string) ([]*challenges.ChallengeSubmission, error) {
	return r.challenges.ChallengeSubmissions(ctx, challengeID)
}

func (r *queryResolver) ChallengeVotes(ctx context.Context, submissionID string) ([]*challenges.ChallengeVote, error) {
	return r.challenges.ChallengeVotes(ctx, submissionID)
}

// mutationResolver - implementa MutationResolver
type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, input users.CreateUserInput) (*users.GraphQLUser, error) {
	return r.users.CreateUser(ctx, input)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input users.UpdateUserInput) (*users.GraphQLUser, error) {
	return r.users.UpdateUser(ctx, id, input)
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	return r.users.DeleteUser(ctx, id)
}

func (r *mutationResolver) CreateChallenge(ctx context.Context, input challenges.CreateChallengeInput) (*challenges.Challenge, error) {
	return r.challenges.CreateChallenge(ctx, input)
}

func (r *mutationResolver) SubmitChallenge(ctx context.Context, input challenges.SubmitChallengeInput) (*challenges.ChallengeSubmission, error) {
	return r.challenges.SubmitChallenge(ctx, input)
}

func (r *mutationResolver) VoteChallenge(ctx context.Context, input challenges.VoteChallengeInput) (*challenges.ChallengeVote, error) {
	return r.challenges.VoteChallenge(ctx, input)
}
