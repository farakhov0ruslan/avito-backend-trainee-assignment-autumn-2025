package main

import (
	"context"
	"fmt"
	"os"

	"avito-backend-trainee-assignment-autumn-2025/internal/config"
	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models"
	"avito-backend-trainee-assignment-autumn-2025/internal/repository/postgres"
	"avito-backend-trainee-assignment-autumn-2025/pkg/database"
	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Init(cfg.App.LogLevel)
	logger.Info("Starting PR Reviewer Assignment Service (Test Mode)")

	// Connect to database
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
		MaxConns: cfg.Database.MaxConns,
		MinConns: cfg.Database.MinConns,
	}

	pool, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	logger.Info("Successfully connected to database")

	// Create repository
	teamRepo := postgres.NewTeamRepository(pool)

	// Test creating a team
	ctx := context.Background()
	testTeam := &models.Team{
		Name:    "Backend-Team",
		Members: []models.User{},
	}

	logger.Info("Creating test team: %s", testTeam.Name)
	err = teamRepo.Create(ctx, testTeam)
	if err != nil {
		logger.Error("Failed to create team: %v", err)
	} else {
		logger.Info("Successfully created team: %s", testTeam.Name)
	}

	// Test retrieving the team
	logger.Info("Retrieving team: %s", testTeam.Name)
	retrievedTeam, err := teamRepo.GetByName(ctx, testTeam.Name)
	if err != nil {
		logger.Error("Failed to retrieve team: %v", err)
		os.Exit(1)
	}

	logger.Info("Successfully retrieved team: %s", retrievedTeam.Name)
	logger.Info("Team has %d members", len(retrievedTeam.Members))

	// Test checking if team exists
	logger.Info("Checking if team exists: %s", testTeam.Name)
	exists, err := teamRepo.Exists(ctx, testTeam.Name)
	if err != nil {
		logger.Error("Failed to check team existence: %v", err)
		os.Exit(1)
	}

	if exists {
		logger.Info("Team exists: %s", testTeam.Name)
	} else {
		logger.Error("Team should exist but doesn't: %s", testTeam.Name)
		os.Exit(1)
	}

	// Test creating duplicate team (should fail)
	logger.Info("Attempting to create duplicate team (should fail)")
	err = teamRepo.Create(ctx, testTeam)
	if err != nil {
		logger.Info("Expected error occurred: %v", err)
	} else {
		logger.Error("Creating duplicate team should have failed but didn't")
		os.Exit(1)
	}

	logger.Info("All tests passed successfully!")
}
