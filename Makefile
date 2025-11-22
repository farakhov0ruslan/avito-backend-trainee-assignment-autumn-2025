.PHONY: help install-goose migrate-up migrate-down migrate-create migrate-status build run test lint docker-up docker-down clean

# Load environment variables from .env file
include .env
export

# Database connection string for goose
DB_DSN := "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"

# Help command
help:
	@echo "Available commands:"
	@echo "  make install-goose       Install goose migration tool"
	@echo "  make migrate-up          Apply all migrations"
	@echo "  make migrate-down        Rollback last migration"
	@echo "  make migrate-create NAME=<name>  Create new migration"
	@echo "  make migrate-status      Show migration status"
	@echo "  make build               Build the application"
	@echo "  make run                 Run the application"
	@echo "  make test                Run all tests"
	@echo "  make test-coverage       Run tests with coverage"
	@echo "  make lint                Run linter"
	@echo "  make fmt                 Format code"
	@echo "  make docker-up           Start services with docker-compose"
	@echo "  make docker-down         Stop services with docker-compose"
	@echo "  make docker-logs         Show docker logs"
	@echo "  make clean               Clean build artifacts"

# Install goose
install-goose:
	@echo "Installing goose..."
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "Goose installed successfully!"

# Apply all migrations
migrate-up:
	@echo "Applying migrations..."
	@goose -dir migrations postgres $(DB_DSN) up
	@echo "Migrations applied successfully!"

# Rollback last migration
migrate-down:
	@echo "Rolling back last migration..."
	@goose -dir migrations postgres $(DB_DSN) down
	@echo "Migration rolled back successfully!"

# Create new migration
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)..."
	@goose -dir migrations create $(NAME) sql
	@echo "Migration created successfully!"

# Show migration status
migrate-status:
	@goose -dir migrations postgres $(DB_DSN) status

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/api cmd/api/main.go
	@echo "Build completed: bin/api"

# Run the application
run:
	@echo "Running application..."
	@go run cmd/api/main.go

# Run all tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

# Start services with docker-compose
docker-up:
	@echo "Starting services with docker-compose..."
	@docker-compose up --build

# Stop services with docker-compose
docker-down:
	@echo "Stopping services..."
	@docker-compose down

# Show docker logs
docker-logs:
	@docker-compose logs -f

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean completed!"

# Install all development tools
install-tools: install-goose
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "All tools installed successfully!"
