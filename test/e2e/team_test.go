package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// TestTeamCreate tests POST /team/add endpoint
func TestTeamCreate(t *testing.T) {
	t.Run("Success - Create team with members", func(t *testing.T) {
		teamName := fmt.Sprintf("test-team-%d", generateID())

		requestBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{
					"user_id":   fmt.Sprintf("user-%d-1", generateID()),
					"username":  "Alice",
					"is_active": true,
				},
				{
					"user_id":   fmt.Sprintf("user-%d-2", generateID()),
					"username":  "Bob",
					"is_active": true,
				},
				{
					"user_id":   fmt.Sprintf("user-%d-3", generateID()),
					"username":  "Charlie",
					"is_active": false,
				},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", requestBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		defer resp.Body.Close()

		assertStatusCode(t, resp, http.StatusCreated)

		var response struct {
			Team struct {
				TeamName string `json:"team_name"`
				Members  []struct {
					UserID   string `json:"user_id"`
					Username string `json:"username"`
					IsActive bool   `json:"is_active"`
				} `json:"members"`
			} `json:"team"`
		}

		if err := parseResponse(resp, &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response.Team.TeamName != teamName {
			t.Errorf("Expected team_name %s, got %s", teamName, response.Team.TeamName)
		}

		if len(response.Team.Members) != 3 {
			t.Errorf("Expected 3 members, got %d", len(response.Team.Members))
		}
	})

	t.Run("Error - Team already exists", func(t *testing.T) {
		teamName := fmt.Sprintf("duplicate-team-%d", generateID())

		requestBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{
					"user_id":   fmt.Sprintf("user-%d-1", generateID()),
					"username":  "Alice",
					"is_active": true,
				},
			},
		}

		// Create team first time
		resp, err := doRequest(http.MethodPost, "/team/add", requestBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		// Try to create same team again
		resp, err = doRequest(http.MethodPost, "/team/add", requestBody)
		if err != nil {
			t.Fatalf("Failed to make second request: %v", err)
		}

		assertStatusCode(t, resp, http.StatusBadRequest)
		assertErrorCode(t, resp, "TEAM_EXISTS")
	})

	t.Run("Success - Create team with no members", func(t *testing.T) {
		teamName := fmt.Sprintf("empty-team-%d", generateID())

		requestBody := map[string]interface{}{
			"team_name": teamName,
			"members":   []map[string]interface{}{},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", requestBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		defer resp.Body.Close()

		assertStatusCode(t, resp, http.StatusCreated)
	})
}

// TestTeamGet tests GET /team/get endpoint
func TestTeamGet(t *testing.T) {
	t.Run("Success - Get existing team", func(t *testing.T) {
		// First create a team
		teamName := fmt.Sprintf("get-team-%d", generateID())
		userID1 := fmt.Sprintf("user-%d-1", generateID())
		userID2 := fmt.Sprintf("user-%d-2", generateID())

		createBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": userID1, "username": "Alice", "is_active": true},
				{"user_id": userID2, "username": "Bob", "is_active": false},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", createBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		// Now get the team
		resp, err = doGet("/team/get", map[string]string{"team_name": teamName})
		if err != nil {
			t.Fatalf("Failed to get team: %v", err)
		}
		defer resp.Body.Close()

		assertStatusCode(t, resp, http.StatusOK)

		var response struct {
			TeamName string `json:"team_name"`
			Members  []struct {
				UserID   string `json:"user_id"`
				Username string `json:"username"`
				IsActive bool   `json:"is_active"`
			} `json:"members"`
		}

		if err := parseResponse(resp, &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response.TeamName != teamName {
			t.Errorf("Expected team_name %s, got %s", teamName, response.TeamName)
		}

		if len(response.Members) != 2 {
			t.Errorf("Expected 2 members, got %d", len(response.Members))
		}
	})

	t.Run("Error - Team not found", func(t *testing.T) {
		resp, err := doGet("/team/get", map[string]string{"team_name": "nonexistent-team"})
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		assertStatusCode(t, resp, http.StatusNotFound)
		assertErrorCode(t, resp, "NOT_FOUND")
	})
}

// generateID generates a unique ID based on current timestamp (nanoseconds)
func generateID() int64 {
	return time.Now().UnixNano()
}
