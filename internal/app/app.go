package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	"avito-backend-trainee-assignment-autumn-2025/internal/config"
	"avito-backend-trainee-assignment-autumn-2025/internal/handler"
	"avito-backend-trainee-assignment-autumn-2025/internal/middleware"
	"avito-backend-trainee-assignment-autumn-2025/internal/repository"
	"avito-backend-trainee-assignment-autumn-2025/internal/repository/postgres"
	"avito-backend-trainee-assignment-autumn-2025/internal/service"
	"avito-backend-trainee-assignment-autumn-2025/pkg/database"
	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

// App represents the application with all its dependencies
type App struct {
	config *config.Config
	db     *pgxpool.Pool
	router *mux.Router
	server *http.Server
}

// NewApp creates and initializes a new application instance
func NewApp(cfg *config.Config) (*App, error) {
	// Initialize logger
	logger.Init(cfg.App.LogLevel)
	logger.Info("Initializing PR Reviewer Assignment Service...")

	// Initialize database connection
	dbConfig := database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DBName:          cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxConns:        cfg.Database.MaxConns,
		MinConns:        cfg.Database.MinConns,
		MaxConnLifetime: cfg.Database.MaxConnLifetime,
		MaxConnIdleTime: cfg.Database.MaxConnIdleTime,
	}

	pool, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize transaction manager
	txManager := repository.NewPgxTransactionManager(pool)

	// Initialize repositories
	teamRepo := postgres.NewTeamRepository(pool)
	userRepo := postgres.NewUserRepository(pool)
	prRepo := postgres.NewPRRepository(pool)

	logger.Info("Repositories initialized")

	// Initialize services
	teamService := service.NewTeamService(teamRepo, userRepo, txManager)
	userService := service.NewUserService(userRepo, prRepo)
	prService := service.NewPRService(prRepo, userRepo, teamRepo, txManager)

	logger.Info("Services initialized")

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	teamHandler := handler.NewTeamHandler(teamService)
	userHandler := handler.NewUserHandler(userService)
	prHandler := handler.NewPRHandler(prService)

	logger.Info("Handlers initialized")

	// Initialize router
	router := NewRouter(healthHandler, teamHandler, userHandler, prHandler)

	logger.Info("Router initialized with all endpoints")

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		config: cfg,
		db:     pool,
		router: router,
		server: server,
	}, nil
}

// NewRouter creates and configures the HTTP router with all endpoints and middleware
func NewRouter(
	healthHandler *handler.HealthHandler,
	teamHandler *handler.TeamHandler,
	userHandler *handler.UserHandler,
	prHandler *handler.PRHandler,
) *mux.Router {
	router := mux.NewRouter()

	// Apply middleware (order matters: Recovery -> Logger)
	router.Use(middleware.Recovery)
	router.Use(middleware.Logger)

	// Health endpoint
	router.HandleFunc("/health", healthHandler.Check).Methods(http.MethodGet)

	// Team endpoints
	router.HandleFunc("/team/add", teamHandler.CreateTeam).Methods(http.MethodPost)
	router.HandleFunc("/team/get", teamHandler.GetTeam).Methods(http.MethodGet)

	// User endpoints
	router.HandleFunc("/users/setIsActive", userHandler.SetIsActive).Methods(http.MethodPost)
	router.HandleFunc("/users/getReview", userHandler.GetUserReviews).Methods(http.MethodGet)

	// Pull Request endpoints
	router.HandleFunc("/pullRequest/create", prHandler.CreatePR).Methods(http.MethodPost)
	router.HandleFunc("/pullRequest/merge", prHandler.MergePR).Methods(http.MethodPost)
	router.HandleFunc("/pullRequest/reassign", prHandler.ReassignReviewer).Methods(http.MethodPost)

	return router
}

// Run starts the HTTP server and handles graceful shutdown
func (a *App) Run() error {
	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Channel to listen for server errors
	serverErrors := make(chan error, 1)

	// Start HTTP server in a goroutine
	go func() {
		logger.Info("Starting HTTP server on port %s", a.config.Server.Port)
		logger.Info("Server is ready to handle requests at http://localhost:%s", a.config.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Wait for interrupt signal or server error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-stop:
		logger.Info("Received signal: %v. Starting graceful shutdown...", sig)
	}

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger.Info("Shutting down HTTP server...")
	if err := a.server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: %v", err)
		return err
	}

	logger.Info("HTTP server stopped")

	// Close database connection
	database.Close(a.db)

	logger.Info("Application shutdown complete")
	return nil
}
