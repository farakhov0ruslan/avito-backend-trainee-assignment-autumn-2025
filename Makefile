.PHONY: help install-goose migrate-up migrate-down migrate-create migrate-status build run test lint docker-up docker-down clean e2e-setup e2e-run-api e2e-test e2e-teardown e2e

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
	@echo "  make e2e                 Show E2E test instructions"
	@echo "  make e2e-setup           Start E2E test database"
	@echo "  make e2e-run-api         Start API server for E2E tests (run in separate terminal)"
	@echo "  make e2e-test            Run E2E tests (requires e2e-setup and e2e-run-api)"
	@echo "  make e2e-teardown        Stop E2E test environment"
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

# E2E Testing Commands

# Start E2E test environment (database only)
e2e-setup:
	@echo "Starting E2E test environment..."
	@if command -v docker-compose &> /dev/null; then \
		docker-compose -f docker-compose.e2e.yml up -d; \
	elif command -v docker &> /dev/null && docker compose version &> /dev/null; then \
		docker compose -f docker-compose.e2e.yml up -d; \
	else \
		echo "Error: docker-compose or docker compose not found"; \
		exit 1; \
	fi
	@echo "Waiting for database to be ready..."
	@sleep 8
	@echo "Running migrations on test database..."
	@goose -dir migrations postgres "postgres://postgres:postgres_test@localhost:5455/pr_reviewer_test_db?sslmode=disable" up
	@echo "E2E database ready!"
	@echo ""
	@echo "To start API server manually, run in another terminal:"
	@echo "  SERVER_PORT=8082 DB_HOST=localhost DB_PORT=5455 DB_USER=postgres DB_PASSWORD=postgres_test DB_NAME=pr_reviewer_test_db DB_SSLMODE=disable go run cmd/api/main.go"
	@echo ""
	@echo "Or use: make e2e-run-api"

# Run API server for E2E tests (run in separate terminal)
e2e-run-api:
	@echo "Starting API server for E2E tests..."
	@SERVER_PORT=8082 \
		DB_HOST=localhost \
		DB_PORT=5455 \
		DB_USER=postgres \
		DB_PASSWORD=postgres_test \
		DB_NAME=pr_reviewer_test_db \
		DB_SSLMODE=disable \
		go run cmd/api/main.go

# Run E2E tests (requires database and API to be running)
e2e-test:
	@echo "Running E2E tests..."
	@echo "Make sure API server is running on port 8082"
	@SERVER_PORT=8082 go test -v ./test/e2e/...

# Stop E2E test environment
e2e-teardown:
	@echo "Stopping E2E test environment..."
	@if command -v docker-compose &> /dev/null; then \
		docker-compose -f docker-compose.e2e.yml down -v; \
	elif command -v docker &> /dev/null && docker compose version &> /dev/null; then \
		docker compose -f docker-compose.e2e.yml down -v; \
	fi
	@echo "E2E environment stopped!"

# Run complete E2E test suite (interactive - requires manual API start)
e2e:
	@echo "==================================================================="
	@echo "E2E Test Suite - Manual Mode"
	@echo "==================================================================="
	@echo ""
	@echo "Step 1: Setting up test database..."
	@$(MAKE) e2e-setup
	@echo ""
	@echo "Step 2: Start API server in another terminal with:"
	@echo "    make e2e-run-api"
	@echo ""
	@echo "Step 3: Run tests with:"
	@echo "    make e2e-test"
	@echo ""
	@echo "Step 4: Clean up with:"
	@echo "    make e2e-teardown"
	@echo ""
	@echo "==================================================================="
