# Testing Guide - PR Reviewer Assignment Service

This guide explains how to test the application after implementing Phase 9.

## Prerequisites

- Docker and Docker Compose installed
- Go 1.24+ installed
- `jq` installed for JSON parsing (optional, for better output in test script)

## Step 1: Start the Database

Start PostgreSQL using docker-compose:

```bash
docker-compose up -d
```

Wait for the database to be ready (check with `docker-compose logs -f postgres`).

## Step 2: Apply Migrations

Apply all database migrations:

```bash
make migrate-up
```

You should see output indicating that migrations were applied successfully.

## Step 3: Build the Application

Build the application binary:

```bash
make build
```

This creates `bin/api` executable.

## Step 4: Run the Application

Start the application in one terminal:

```bash
make run
```

Or run the built binary:

```bash
./bin/api
```

You should see output like:
```
INFO: Initializing PR Reviewer Assignment Service...
INFO: Successfully connected to PostgreSQL database: pr_reviewer_db
INFO: Repositories initialized
INFO: Services initialized
INFO: Handlers initialized
INFO: Router initialized with all endpoints
INFO: Starting HTTP server on port 8080
INFO: Server is ready to handle requests at http://localhost:8080
```

## Step 5: Test the API

In another terminal, run the test script:

```bash
./test-api.sh
```

This script will:
1. Check health endpoint
2. Create a team with 3 members
3. Retrieve the team
4. Create a pull request (with automatic reviewer assignment)
5. Get user reviews
6. Deactivate a user
7. Reassign a reviewer
8. Merge the PR
9. Test idempotency of merge
10. Test that reassignment fails on merged PR

## Step 6: Manual Testing with curl

You can also test endpoints manually:

### Health Check
```bash
curl http://localhost:8080/health
```

### Create Team
```bash
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "frontend",
    "members": [
      {"user_id": "u10", "username": "David", "is_active": true},
      {"user_id": "u11", "username": "Eve", "is_active": true}
    ]
  }'
```

### Get Team
```bash
curl "http://localhost:8080/team/get?team_name=frontend"
```

### Create Pull Request
```bash
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-2001",
    "pull_request_name": "Fix bug in UI",
    "author_id": "u10"
  }'
```

### Get User Reviews
```bash
curl "http://localhost:8080/users/getReview?user_id=u11"
```

### Set User Active Status
```bash
curl -X POST http://localhost:8080/users/setIsActive \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "u11",
    "is_active": false
  }'
```

### Merge PR
```bash
curl -X POST http://localhost:8080/pullRequest/merge \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-2001"
  }'
```

### Reassign Reviewer
```bash
curl -X POST http://localhost:8080/pullRequest/reassign \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-2001",
    "old_user_id": "u11"
  }'
```

## Step 7: Check Logs

Watch the application logs to see request handling:

```
INFO: HTTP POST /team/add - Status: 201 - Duration: 45ms
INFO: HTTP GET /team/get - Status: 200 - Duration: 12ms
INFO: HTTP POST /pullRequest/create - Status: 201 - Duration: 38ms
```

## Step 8: Stop Everything

Stop the application with Ctrl+C (it will shutdown gracefully).

Stop the database:

```bash
docker-compose down
```

## Expected Behavior

### Success Cases
- ✅ Health check returns 200
- ✅ Team creation returns 201 with team data
- ✅ PR creation assigns up to 2 active reviewers from author's team
- ✅ Merge is idempotent (can call multiple times)
- ✅ All responses follow OpenAPI specification

### Error Cases
- ❌ Creating duplicate team returns 409 with TEAM_EXISTS code
- ❌ Creating PR with non-existent author returns 404
- ❌ Reassigning on merged PR returns 409 with PR_MERGED code
- ❌ Reassigning non-assigned reviewer returns 409 with NOT_ASSIGNED code

## Troubleshooting

### Database Connection Failed
- Check if PostgreSQL is running: `docker-compose ps`
- Check database logs: `docker-compose logs postgres`
- Verify .env file has correct DB_PORT (should match docker-compose port)

### Migrations Failed
- Ensure database is running
- Check migration files in `migrations/` directory
- Run `make migrate-status` to see current state

### Port Already in Use
- Change SERVER_PORT in .env file
- Kill process using port 8080: `lsof -ti:8080 | xargs kill`

## Performance Testing

The application should meet these requirements:
- Response time < 300ms (SLI)
- Success rate > 99.9% (SLI)
- Handle up to 5 RPS
- Support up to 20 teams and 200 users

Monitor response times in logs to verify performance.
