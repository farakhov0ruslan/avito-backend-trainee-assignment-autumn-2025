package main

import (
	"fmt"
	"os"

	"avito-backend-trainee-assignment-autumn-2025/internal/app"
	"avito-backend-trainee-assignment-autumn-2025/internal/config"
	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Create and initialize application
	application, err := app.NewApp(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	// Run application (blocks until shutdown signal)
	if err := application.Run(); err != nil {
		logger.Error("Application error: %v", err)
		os.Exit(1)
	}
}
