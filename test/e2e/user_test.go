package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// TestUserSetActive tests POST /users/setIsActive endpoint
func TestUserSetActive(t *testing.T) {
	t.Run("Success - Set user to inactive", func(t *testing.T) {
		// Create a team with a user first
		teamName := fmt.Sprintf("user-team-%d", time.Now().UnixNano())
		userID := fmt.Sprintf("user-%d", time.Now().UnixNano())

		createBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": userID, "username": "TestUser", "is_active": true},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", createBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		// Set user to inactive
		setActiveBody := map[string]interface{}{
			"user_id":   userID,
			"is_active": false,
		}

		resp, err = doRequest(http.MethodPost, "/users/setIsActive", setActiveBody)
		if err != nil {
			t.Fatalf("Failed to set user active: %v", err)
		}
		defer resp.Body.Close()

		assertStatusCode(t, resp, http.StatusOK)

		var response struct {
			User struct {
				UserID   string `json:"user_id"`
				Username string `json:"username"`
				TeamName string `json:"team_name"`
				IsActive bool   `json:"is_active"`
			} `json:"user"`
		}

		if err := parseResponse(resp, &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response.User.UserID != userID {
			t.Errorf("Expected user_id %s, got %s", userID, response.User.UserID)
		}

		if response.User.IsActive {
			t.Error("Expected is_active to be false")
		}

		if response.User.TeamName != teamName {
			t.Errorf("Expected team_name %s, got %s", teamName, response.User.TeamName)
		}
	})

	t.Run("Success - Set user to active", func(t *testing.T) {
		// Create a team with inactive user
		teamName := fmt.Sprintf("user-team-%d", time.Now().UnixNano())
		userID := fmt.Sprintf("user-%d", time.Now().UnixNano())

		createBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": userID, "username": "TestUser", "is_active": false},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", createBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		// Set user to active
		setActiveBody := map[string]interface{}{
			"user_id":   userID,
			"is_active": true,
		}

		resp, err = doRequest(http.MethodPost, "/users/setIsActive", setActiveBody)
		if err != nil {
			t.Fatalf("Failed to set user active: %v", err)
		}
		defer resp.Body.Close()

		assertStatusCode(t, resp, http.StatusOK)

		var response struct {
			User struct {
				IsActive bool `json:"is_active"`
			} `json:"user"`
		}

		if err := parseResponse(resp, &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if !response.User.IsActive {
			t.Error("Expected is_active to be true")
		}
	})

	t.Run("Error - User not found", func(t *testing.T) {
		setActiveBody := map[string]interface{}{
			"user_id":   "nonexistent-user",
			"is_active": true,
		}

		resp, err := doRequest(http.MethodPost, "/users/setIsActive", setActiveBody)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		assertStatusCode(t, resp, http.StatusNotFound)
		assertErrorCode(t, resp, "NOT_FOUND")
	})
}

// TestUserGetReviews tests GET /users/getReview endpoint
func TestUserGetReviews(t *testing.T) {
	t.Run("Success - Get user reviews (no PRs)", func(t *testing.T) {
		// Create a team with a user
		teamName := fmt.Sprintf("review-team-%d", time.Now().UnixNano())
		userID := fmt.Sprintf("user-%d", time.Now().UnixNano())

		createBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": userID, "username": "Reviewer", "is_active": true},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", createBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		// Get user reviews
		resp, err = doGet("/users/getReview", map[string]string{"user_id": userID})
		if err != nil {
			t.Fatalf("Failed to get user reviews: %v", err)
		}
		defer resp.Body.Close()

		assertStatusCode(t, resp, http.StatusOK)

		var response struct {
			UserID       string `json:"user_id"`
			PullRequests []struct {
				PullRequestID   string `json:"pull_request_id"`
				PullRequestName string `json:"pull_request_name"`
				AuthorID        string `json:"author_id"`
				Status          string `json:"status"`
			} `json:"pull_requests"`
		}

		if err := parseResponse(resp, &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response.UserID != userID {
			t.Errorf("Expected user_id %s, got %s", userID, response.UserID)
		}

		if len(response.PullRequests) != 0 {
			t.Errorf("Expected 0 pull requests, got %d", len(response.PullRequests))
		}
	})

	t.Run("Success - Get user reviews (with PRs)", func(t *testing.T) {
		// Create a team with multiple users
		teamName := fmt.Sprintf("review-team-%d", time.Now().UnixNano())
		authorID := fmt.Sprintf("author-%d", time.Now().UnixNano())
		reviewerID := fmt.Sprintf("reviewer-%d", time.Now().UnixNano())
		reviewer2ID := fmt.Sprintf("reviewer2-%d", time.Now().UnixNano())

		createBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": authorID, "username": "Author", "is_active": true},
				{"user_id": reviewerID, "username": "Reviewer", "is_active": true},
				{"user_id": reviewer2ID, "username": "Reviewer2", "is_active": true},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", createBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		// Create PRs
		pr1ID := fmt.Sprintf("pr-%d-1", time.Now().UnixNano())
		pr2ID := fmt.Sprintf("pr-%d-2", time.Now().UnixNano())

		prBody1 := map[string]interface{}{
			"pull_request_id":   pr1ID,
			"pull_request_name": "Feature A",
			"author_id":         authorID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/create", prBody1)
		if err != nil {
			t.Fatalf("Failed to create PR 1: %v", err)
		}
		resp.Body.Close()

		prBody2 := map[string]interface{}{
			"pull_request_id":   pr2ID,
			"pull_request_name": "Feature B",
			"author_id":         authorID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/create", prBody2)
		if err != nil {
			t.Fatalf("Failed to create PR 2: %v", err)
		}
		resp.Body.Close()

		// Get reviewer's reviews
		resp, err = doGet("/users/getReview", map[string]string{"user_id": reviewerID})
		if err != nil {
			t.Fatalf("Failed to get user reviews: %v", err)
		}
		defer resp.Body.Close()

		assertStatusCode(t, resp, http.StatusOK)

		var response struct {
			UserID       string `json:"user_id"`
			PullRequests []struct {
				PullRequestID   string `json:"pull_request_id"`
				PullRequestName string `json:"pull_request_name"`
				AuthorID        string `json:"author_id"`
				Status          string `json:"status"`
			} `json:"pull_requests"`
		}

		if err := parseResponse(resp, &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response.UserID != reviewerID {
			t.Errorf("Expected user_id %s, got %s", reviewerID, response.UserID)
		}

		// User should be assigned to at least some PRs (randomized, so can't guarantee exact count)
		t.Logf("Reviewer %s is assigned to %d PR(s)", reviewerID, len(response.PullRequests))
	})

	t.Run("Error - User not found", func(t *testing.T) {
		resp, err := doGet("/users/getReview", map[string]string{"user_id": "nonexistent-user"})
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		assertStatusCode(t, resp, http.StatusNotFound)
		assertErrorCode(t, resp, "NOT_FOUND")
	})
}
