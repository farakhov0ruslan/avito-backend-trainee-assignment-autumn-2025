package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string

	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// NewPostgresDB creates a new PostgreSQL connection pool
func NewPostgresDB(cfg Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	// Configure connection pool
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Ping database with retries (5 attempts)
	const maxRetries = 5
	const retryDelay = 2 * time.Second

	var pingErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		pingErr = pool.Ping(ctx)
		if pingErr == nil {
			logger.Info("Successfully connected to PostgreSQL database: %s", cfg.DBName)
			return pool, nil
		}

		if attempt < maxRetries {
			logger.Warn("Failed to ping database (attempt %d/%d): %v. Retrying in %v...",
				attempt, maxRetries, pingErr, retryDelay)
			time.Sleep(retryDelay)
		}
	}

	pool.Close()
	return nil, fmt.Errorf("unable to ping database after %d attempts: %w", maxRetries, pingErr)
}

// Close gracefully closes the database connection pool
func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
		logger.Info("Database connection pool closed")
	}
}
