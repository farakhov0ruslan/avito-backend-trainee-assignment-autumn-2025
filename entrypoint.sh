#!/usr/bin/env sh
set -e

echo "[entrypoint] Starting PR Reviewer service..."

# Database configuration with defaults
DB_HOST="${DB_HOST:-postgres}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-pr_reviewer_db}"
DB_SSLMODE="${DB_SSLMODE:-disable}"

# Build connection string
DB_DSN="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

echo "[entrypoint] Waiting for database to be ready..."
# Simple wait loop for database
for i in $(seq 1 30); do
    if goose -dir /migrations postgres "${DB_DSN}" version >/dev/null 2>&1; then
        echo "[entrypoint] Database is ready!"
        break
    fi
    echo "[entrypoint] Waiting for database... ($i/30)"
    sleep 2
done

echo "[entrypoint] Running database migrations..."
goose -dir /migrations postgres "${DB_DSN}" up

if [ $? -eq 0 ]; then
    echo "[entrypoint] Migrations applied successfully!"
else
    echo "[entrypoint] ERROR: Failed to apply migrations"
    exit 1
fi

echo "[entrypoint] Starting API server on port ${SERVER_PORT:-8080}..."
exec /app/api
