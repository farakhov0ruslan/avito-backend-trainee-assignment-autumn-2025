#!/bin/bash

# Script to reset database completely

set -e

# Detect docker compose command (docker-compose or docker compose)
if command -v docker-compose &> /dev/null; then
    DOCKER_COMPOSE="docker-compose"
elif command -v docker &> /dev/null && docker compose version &> /dev/null; then
    DOCKER_COMPOSE="docker compose"
else
    echo "Error: Neither 'docker-compose' nor 'docker compose' found"
    exit 1
fi

echo "🔄 Resetting database..."

# Stop containers if running
echo "Stopping containers..."
$DOCKER_COMPOSE down -v 2>/dev/null || true

# Start database
echo "Starting PostgreSQL..."
$DOCKER_COMPOSE up -d

# Wait for database to be ready
echo "Waiting for database to be ready..."
sleep 10

# Apply migrations
echo "Applying migrations..."
make migrate-up

echo "✅ Database reset complete!"
echo ""
echo "You can now run: make run"
