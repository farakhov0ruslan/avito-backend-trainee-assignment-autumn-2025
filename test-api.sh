#!/bin/bash

# Test script for PR Reviewer Assignment Service
# This script tests all API endpoints to verify they work correctly

set -e  # Exit on error

BASE_URL="http://localhost:8080"
BOLD='\033[1m'
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${BOLD}=== PR Reviewer Assignment Service - API Tests ===${NC}\n"

# Function to print test header
print_test() {
    echo -e "${BOLD}Test: $1${NC}"
}

# Function to print success
print_success() {
    echo -e "${GREEN}✓ $1${NC}\n"
}

# Function to print error
print_error() {
    echo -e "${RED}✗ $1${NC}\n"
}

# Function to print info
print_info() {
    echo -e "${YELLOW}→ $1${NC}"
}

# Check if server is running
print_test "1. Health Check"
print_info "GET /health"
if curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/health" | grep -q "200"; then
    print_success "Server is running"
else
    print_error "Server is not responding"
    exit 1
fi

# Test 1: Create a team with 4 members (need enough for reassign test)
print_test "2. Create Team (Backend Team)"
print_info "POST /team/add"
RESPONSE=$(curl -s -X POST "$BASE_URL/team/add" \
    -H "Content-Type: application/json" \
    -d '{
        "team_name": "backend",
        "members": [
            {"user_id": "u1", "username": "Alice", "is_active": true},
            {"user_id": "u2", "username": "Bob", "is_active": true},
            {"user_id": "u3", "username": "Charlie", "is_active": true},
            {"user_id": "u4", "username": "Dave", "is_active": true}
        ]
    }')
echo "$RESPONSE" | jq '.'
if echo "$RESPONSE" | jq -e '.team.team_name == "backend"' > /dev/null; then
    print_success "Team created successfully"
else
    print_error "Failed to create team"
fi

# Test 2: Get team
print_test "3. Get Team"
print_info "GET /team/get?team_name=backend"
RESPONSE=$(curl -s "$BASE_URL/team/get?team_name=backend")
echo "$RESPONSE" | jq '.'
if echo "$RESPONSE" | jq -e '.team_name == "backend"' > /dev/null; then
    print_success "Team retrieved successfully"
else
    print_error "Failed to get team"
fi

# Test 3: Create PR
print_test "4. Create Pull Request"
print_info "POST /pullRequest/create"
RESPONSE=$(curl -s -X POST "$BASE_URL/pullRequest/create" \
    -H "Content-Type: application/json" \
    -d '{
        "pull_request_id": "pr-1001",
        "pull_request_name": "Add search feature",
        "author_id": "u1"
    }')
echo "$RESPONSE" | jq '.'
if echo "$RESPONSE" | jq -e '.pr.pull_request_id == "pr-1001"' > /dev/null; then
    print_success "PR created successfully with reviewers assigned"
    echo "$RESPONSE" | jq -r '.pr.assigned_reviewers | join(", ")' | xargs -I {} echo "  Assigned reviewers: {}"
else
    print_error "Failed to create PR"
fi

# Test 4: Get user reviews
print_test "5. Get User Reviews"
print_info "GET /users/getReview?user_id=u2"
RESPONSE=$(curl -s "$BASE_URL/users/getReview?user_id=u2")
echo "$RESPONSE" | jq '.'
print_success "User reviews retrieved"

# Test 5: Set user inactive
print_test "6. Set User Inactive"
print_info "POST /users/setIsActive"
RESPONSE=$(curl -s -X POST "$BASE_URL/users/setIsActive" \
    -H "Content-Type: application/json" \
    -d '{
        "user_id": "u3",
        "is_active": false
    }')
echo "$RESPONSE" | jq '.'
if echo "$RESPONSE" | jq -e '.user.is_active == false' > /dev/null; then
    print_success "User deactivated successfully"
else
    print_error "Failed to deactivate user"
fi

# Test 6: Reassign reviewer
print_test "7. Reassign Reviewer"
print_info "POST /pullRequest/reassign"
RESPONSE=$(curl -s -X POST "$BASE_URL/pullRequest/reassign" \
    -H "Content-Type: application/json" \
    -d '{
        "pull_request_id": "pr-1001",
        "old_user_id": "u2"
    }')
echo "$RESPONSE" | jq '.'
if echo "$RESPONSE" | jq -e '.replaced_by' > /dev/null; then
    NEW_REVIEWER=$(echo "$RESPONSE" | jq -r '.replaced_by')
    print_success "Reviewer reassigned successfully to $NEW_REVIEWER"
else
    print_error "Failed to reassign reviewer"
fi

# Test 7: Merge PR
print_test "8. Merge Pull Request"
print_info "POST /pullRequest/merge"
RESPONSE=$(curl -s -X POST "$BASE_URL/pullRequest/merge" \
    -H "Content-Type: application/json" \
    -d '{
        "pull_request_id": "pr-1001"
    }')
echo "$RESPONSE" | jq '.'
if echo "$RESPONSE" | jq -e '.pr.status == "MERGED"' > /dev/null; then
    print_success "PR merged successfully"
else
    print_error "Failed to merge PR"
fi

# Test 8: Try to merge again (idempotency test)
print_test "9. Merge PR Again (Idempotency Test)"
print_info "POST /pullRequest/merge"
RESPONSE=$(curl -s -X POST "$BASE_URL/pullRequest/merge" \
    -H "Content-Type: application/json" \
    -d '{
        "pull_request_id": "pr-1001"
    }')
echo "$RESPONSE" | jq '.'
if echo "$RESPONSE" | jq -e '.pr.status == "MERGED"' > /dev/null; then
    print_success "Idempotency verified - PR still merged"
else
    print_error "Idempotency failed"
fi

# Test 9: Try to reassign reviewer on merged PR (should fail)
print_test "10. Try to Reassign on Merged PR (Should Fail)"
print_info "POST /pullRequest/reassign"
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "$BASE_URL/pullRequest/reassign" \
    -H "Content-Type: application/json" \
    -d '{
        "pull_request_id": "pr-1001",
        "old_user_id": "u2"
    }')
HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d':' -f2)
if [ "$HTTP_STATUS" = "409" ]; then
    print_success "Correctly rejected reassignment on merged PR (409 Conflict)"
else
    print_error "Should have rejected reassignment on merged PR"
fi

echo -e "${BOLD}${GREEN}=== All Tests Completed ===${NC}"
