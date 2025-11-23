package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// TestPRCreate tests POST /pullRequest/create endpoint
func TestPRCreate(t *testing.T) {
	t.Run("Success - Create PR with 2 reviewers", func(t *testing.T) {
		// Create team with 4 active users (1 author + 3 potential reviewers)
		teamName := fmt.Sprintf("pr-team-%d", time.Now().UnixNano())
		authorID := fmt.Sprintf("author-%d", time.Now().UnixNano())
		user1ID := fmt.Sprintf("user1-%d", time.Now().UnixNano())
		user2ID := fmt.Sprintf("user2-%d", time.Now().UnixNano())
		user3ID := fmt.Sprintf("user3-%d", time.Now().UnixNano())

		createTeamBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": authorID, "username": "Author", "is_active": true},
				{"user_id": user1ID, "username": "User1", "is_active": true},
				{"user_id": user2ID, "username": "User2", "is_active": true},
				{"user_id": user3ID, "username": "User3", "is_active": true},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", createTeamBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		// Create PR
		prID := fmt.Sprintf("pr-%d", time.Now().UnixNano())
		createPRBody := map[string]interface{}{
			"pull_request_id":   prID,
			"pull_request_name": "Add feature",
			"author_id":         authorID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/create", createPRBody)
		if err != nil {
			t.Fatalf("Failed to create PR: %v", err)
		}
		defer resp.Body.Close()

		assertStatusCode(t, resp, http.StatusCreated)

		var response struct {
			PR struct {
				PullRequestID     string   `json:"pull_request_id"`
				PullRequestName   string   `json:"pull_request_name"`
				AuthorID          string   `json:"author_id"`
				Status            string   `json:"status"`
				AssignedReviewers []string `json:"assigned_reviewers"`
			} `json:"pr"`
		}

		if err := parseResponse(resp, &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response.PR.PullRequestID != prID {
			t.Errorf("Expected pull_request_id %s, got %s", prID, response.PR.PullRequestID)
		}

		if response.PR.Status != "OPEN" {
			t.Errorf("Expected status OPEN, got %s", response.PR.Status)
		}

		if len(response.PR.AssignedReviewers) != 2 {
			t.Errorf("Expected 2 reviewers, got %d", len(response.PR.AssignedReviewers))
		}

		// Verify author is not in reviewers
		for _, reviewerID := range response.PR.AssignedReviewers {
			if reviewerID == authorID {
				t.Error("Author should not be assigned as reviewer")
			}
		}
	})

	t.Run("Success - Create PR with 1 reviewer (only 2 active users)", func(t *testing.T) {
		teamName := fmt.Sprintf("small-team-%d", time.Now().UnixNano())
		authorID := fmt.Sprintf("author-%d", time.Now().UnixNano())
		reviewerID := fmt.Sprintf("reviewer-%d", time.Now().UnixNano())

		createTeamBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": authorID, "username": "Author", "is_active": true},
				{"user_id": reviewerID, "username": "Reviewer", "is_active": true},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", createTeamBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		// Create PR
		prID := fmt.Sprintf("pr-%d", time.Now().UnixNano())
		createPRBody := map[string]interface{}{
			"pull_request_id":   prID,
			"pull_request_name": "Fix bug",
			"author_id":         authorID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/create", createPRBody)
		if err != nil {
			t.Fatalf("Failed to create PR: %v", err)
		}
		defer resp.Body.Close()

		assertStatusCode(t, resp, http.StatusCreated)

		var response struct {
			PR struct {
				AssignedReviewers []string `json:"assigned_reviewers"`
			} `json:"pr"`
		}

		if err := parseResponse(resp, &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if len(response.PR.AssignedReviewers) != 1 {
			t.Errorf("Expected 1 reviewer, got %d", len(response.PR.AssignedReviewers))
		}

		if len(response.PR.AssignedReviewers) > 0 && response.PR.AssignedReviewers[0] != reviewerID {
			t.Errorf("Expected reviewer %s, got %s", reviewerID, response.PR.AssignedReviewers[0])
		}
	})

	t.Run("Success - Create PR with 0 reviewers (only author is active)", func(t *testing.T) {
		teamName := fmt.Sprintf("solo-team-%d", time.Now().UnixNano())
		authorID := fmt.Sprintf("author-%d", time.Now().UnixNano())
		inactiveID := fmt.Sprintf("inactive-%d", time.Now().UnixNano())

		createTeamBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": authorID, "username": "Author", "is_active": true},
				{"user_id": inactiveID, "username": "Inactive", "is_active": false},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", createTeamBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		// Create PR
		prID := fmt.Sprintf("pr-%d", time.Now().UnixNano())
		createPRBody := map[string]interface{}{
			"pull_request_id":   prID,
			"pull_request_name": "Solo work",
			"author_id":         authorID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/create", createPRBody)
		if err != nil {
			t.Fatalf("Failed to create PR: %v", err)
		}
		defer resp.Body.Close()

		assertStatusCode(t, resp, http.StatusCreated)

		var response struct {
			PR struct {
				AssignedReviewers []string `json:"assigned_reviewers"`
			} `json:"pr"`
		}

		if err := parseResponse(resp, &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if len(response.PR.AssignedReviewers) != 0 {
			t.Errorf("Expected 0 reviewers, got %d", len(response.PR.AssignedReviewers))
		}
	})

	t.Run("Error - PR already exists", func(t *testing.T) {
		teamName := fmt.Sprintf("dup-team-%d", time.Now().UnixNano())
		authorID := fmt.Sprintf("author-%d", time.Now().UnixNano())

		createTeamBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": authorID, "username": "Author", "is_active": true},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", createTeamBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		// Create PR first time
		prID := fmt.Sprintf("pr-%d", time.Now().UnixNano())
		createPRBody := map[string]interface{}{
			"pull_request_id":   prID,
			"pull_request_name": "Feature",
			"author_id":         authorID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/create", createPRBody)
		if err != nil {
			t.Fatalf("Failed to create PR: %v", err)
		}
		resp.Body.Close()

		// Try to create same PR again
		resp, err = doRequest(http.MethodPost, "/pullRequest/create", createPRBody)
		if err != nil {
			t.Fatalf("Failed to make second request: %v", err)
		}

		assertStatusCode(t, resp, http.StatusConflict)
		assertErrorCode(t, resp, "PR_EXISTS")
	})

	t.Run("Error - Author not found", func(t *testing.T) {
		prID := fmt.Sprintf("pr-%d", time.Now().UnixNano())
		createPRBody := map[string]interface{}{
			"pull_request_id":   prID,
			"pull_request_name": "Feature",
			"author_id":         "nonexistent-author",
		}

		resp, err := doRequest(http.MethodPost, "/pullRequest/create", createPRBody)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		assertStatusCode(t, resp, http.StatusNotFound)
		assertErrorCode(t, resp, "NOT_FOUND")
	})
}

// TestPRMerge tests POST /pullRequest/merge endpoint
func TestPRMerge(t *testing.T) {
	t.Run("Success - Merge PR", func(t *testing.T) {
		// Create team and PR
		teamName := fmt.Sprintf("merge-team-%d", time.Now().UnixNano())
		authorID := fmt.Sprintf("author-%d", time.Now().UnixNano())
		prID := fmt.Sprintf("pr-%d", time.Now().UnixNano())

		createTeamBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": authorID, "username": "Author", "is_active": true},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", createTeamBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		createPRBody := map[string]interface{}{
			"pull_request_id":   prID,
			"pull_request_name": "Feature",
			"author_id":         authorID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/create", createPRBody)
		if err != nil {
			t.Fatalf("Failed to create PR: %v", err)
		}
		resp.Body.Close()

		// Merge PR
		mergeBody := map[string]interface{}{
			"pull_request_id": prID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/merge", mergeBody)
		if err != nil {
			t.Fatalf("Failed to merge PR: %v", err)
		}
		defer resp.Body.Close()

		assertStatusCode(t, resp, http.StatusOK)

		var response struct {
			PR struct {
				PullRequestID string `json:"pull_request_id"`
				Status        string `json:"status"`
				MergedAt      string `json:"mergedAt"`
			} `json:"pr"`
		}

		if err := parseResponse(resp, &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response.PR.Status != "MERGED" {
			t.Errorf("Expected status MERGED, got %s", response.PR.Status)
		}

		if response.PR.MergedAt == "" {
			t.Error("Expected mergedAt to be set")
		}
	})

	t.Run("Success - Merge PR idempotency", func(t *testing.T) {
		// Create team and PR
		teamName := fmt.Sprintf("idempotent-team-%d", time.Now().UnixNano())
		authorID := fmt.Sprintf("author-%d", time.Now().UnixNano())
		prID := fmt.Sprintf("pr-%d", time.Now().UnixNano())

		createTeamBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": authorID, "username": "Author", "is_active": true},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", createTeamBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		createPRBody := map[string]interface{}{
			"pull_request_id":   prID,
			"pull_request_name": "Feature",
			"author_id":         authorID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/create", createPRBody)
		if err != nil {
			t.Fatalf("Failed to create PR: %v", err)
		}
		resp.Body.Close()

		// Merge PR first time
		mergeBody := map[string]interface{}{
			"pull_request_id": prID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/merge", mergeBody)
		if err != nil {
			t.Fatalf("Failed to merge PR first time: %v", err)
		}

		var firstResponse struct {
			PR struct {
				MergedAt string `json:"mergedAt"`
			} `json:"pr"`
		}
		if err := parseResponse(resp, &firstResponse); err != nil {
			t.Fatalf("Failed to parse first response: %v", err)
		}

		// Merge PR second time (should be idempotent)
		resp, err = doRequest(http.MethodPost, "/pullRequest/merge", mergeBody)
		if err != nil {
			t.Fatalf("Failed to merge PR second time: %v", err)
		}
		defer resp.Body.Close()

		assertStatusCode(t, resp, http.StatusOK)

		var secondResponse struct {
			PR struct {
				Status   string `json:"status"`
				MergedAt string `json:"mergedAt"`
			} `json:"pr"`
		}

		if err := parseResponse(resp, &secondResponse); err != nil {
			t.Fatalf("Failed to parse second response: %v", err)
		}

		if secondResponse.PR.Status != "MERGED" {
			t.Errorf("Expected status MERGED, got %s", secondResponse.PR.Status)
		}

		// mergedAt should be the same
		if firstResponse.PR.MergedAt != secondResponse.PR.MergedAt {
			t.Errorf("mergedAt changed between calls: %s vs %s",
				firstResponse.PR.MergedAt, secondResponse.PR.MergedAt)
		}
	})

	t.Run("Error - PR not found", func(t *testing.T) {
		mergeBody := map[string]interface{}{
			"pull_request_id": "nonexistent-pr",
		}

		resp, err := doRequest(http.MethodPost, "/pullRequest/merge", mergeBody)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		assertStatusCode(t, resp, http.StatusNotFound)
		assertErrorCode(t, resp, "NOT_FOUND")
	})
}

// TestPRReassign tests POST /pullRequest/reassign endpoint
func TestPRReassign(t *testing.T) {
	t.Run("Success - Reassign reviewer", func(t *testing.T) {
		// Create team with 4 active users
		teamName := fmt.Sprintf("reassign-team-%d", time.Now().UnixNano())
		authorID := fmt.Sprintf("author-%d", time.Now().UnixNano())
		user1ID := fmt.Sprintf("user1-%d", time.Now().UnixNano())
		user2ID := fmt.Sprintf("user2-%d", time.Now().UnixNano())
		user3ID := fmt.Sprintf("user3-%d", time.Now().UnixNano())

		createTeamBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": authorID, "username": "Author", "is_active": true},
				{"user_id": user1ID, "username": "User1", "is_active": true},
				{"user_id": user2ID, "username": "User2", "is_active": true},
				{"user_id": user3ID, "username": "User3", "is_active": true},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", createTeamBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		// Create PR (will assign 2 reviewers randomly)
		prID := fmt.Sprintf("pr-%d", time.Now().UnixNano())
		createPRBody := map[string]interface{}{
			"pull_request_id":   prID,
			"pull_request_name": "Feature",
			"author_id":         authorID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/create", createPRBody)
		if err != nil {
			t.Fatalf("Failed to create PR: %v", err)
		}

		var createResponse struct {
			PR struct {
				AssignedReviewers []string `json:"assigned_reviewers"`
			} `json:"pr"`
		}

		if err := parseResponse(resp, &createResponse); err != nil {
			t.Fatalf("Failed to parse create response: %v", err)
		}

		if len(createResponse.PR.AssignedReviewers) == 0 {
			t.Skip("No reviewers assigned, skipping reassign test")
		}

		// Get one of the assigned reviewers
		oldReviewerID := createResponse.PR.AssignedReviewers[0]

		// Reassign that reviewer
		reassignBody := map[string]interface{}{
			"pull_request_id": prID,
			"old_user_id":     oldReviewerID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/reassign", reassignBody)
		if err != nil {
			t.Fatalf("Failed to reassign reviewer: %v", err)
		}
		defer resp.Body.Close()

		assertStatusCode(t, resp, http.StatusOK)

		var reassignResponse struct {
			PR struct {
				AssignedReviewers []string `json:"assigned_reviewers"`
			} `json:"pr"`
			ReplacedBy string `json:"replaced_by"`
		}

		if err := parseResponse(resp, &reassignResponse); err != nil {
			t.Fatalf("Failed to parse reassign response: %v", err)
		}

		// Verify old reviewer is no longer in the list
		for _, reviewerID := range reassignResponse.PR.AssignedReviewers {
			if reviewerID == oldReviewerID {
				t.Errorf("Old reviewer %s should not be in assigned_reviewers", oldReviewerID)
			}
		}

		// Verify replaced_by is set
		if reassignResponse.ReplacedBy == "" {
			t.Error("Expected replaced_by to be set")
		}

		// Verify new reviewer is in the list
		found := false
		for _, reviewerID := range reassignResponse.PR.AssignedReviewers {
			if reviewerID == reassignResponse.ReplacedBy {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("New reviewer %s should be in assigned_reviewers", reassignResponse.ReplacedBy)
		}
	})

	t.Run("Error - Reviewer not assigned", func(t *testing.T) {
		// Create team with 4 users so that 2 are assigned as reviewers and 2 are not
		teamName := fmt.Sprintf("notassigned-team-%d", time.Now().UnixNano())
		authorID := fmt.Sprintf("author-%d", time.Now().UnixNano())
		reviewer1ID := fmt.Sprintf("reviewer1-%d", time.Now().UnixNano())
		reviewer2ID := fmt.Sprintf("reviewer2-%d", time.Now().UnixNano())
		reviewer3ID := fmt.Sprintf("reviewer3-%d", time.Now().UnixNano())

		allReviewerIDs := []string{reviewer1ID, reviewer2ID, reviewer3ID}

		createTeamBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": authorID, "username": "Author", "is_active": true},
				{"user_id": reviewer1ID, "username": "Reviewer1", "is_active": true},
				{"user_id": reviewer2ID, "username": "Reviewer2", "is_active": true},
				{"user_id": reviewer3ID, "username": "Reviewer3", "is_active": true},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", createTeamBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		// Create PR
		prID := fmt.Sprintf("pr-%d", time.Now().UnixNano())
		createPRBody := map[string]interface{}{
			"pull_request_id":   prID,
			"pull_request_name": "Feature",
			"author_id":         authorID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/create", createPRBody)
		if err != nil {
			t.Fatalf("Failed to create PR: %v", err)
		}

		// Parse response to see who was assigned
		var createResponse struct {
			PR struct {
				AssignedReviewers []string `json:"assigned_reviewers"`
			} `json:"pr"`
		}
		if err := parseResponse(resp, &createResponse); err != nil {
			t.Fatalf("Failed to parse create response: %v", err)
		}

		// Find a user who was NOT assigned
		var notAssignedID string
		for _, reviewerID := range allReviewerIDs {
			isAssigned := false
			for _, assignedID := range createResponse.PR.AssignedReviewers {
				if reviewerID == assignedID {
					isAssigned = true
					break
				}
			}
			if !isAssigned {
				notAssignedID = reviewerID
				break
			}
		}

		if notAssignedID == "" {
			t.Fatal("All reviewers were assigned, expected at least one to not be assigned")
		}

		// Try to reassign a user who is not assigned
		reassignBody := map[string]interface{}{
			"pull_request_id": prID,
			"old_user_id":     notAssignedID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/reassign", reassignBody)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		assertStatusCode(t, resp, http.StatusConflict)
		assertErrorCode(t, resp, "NOT_ASSIGNED")
	})

	t.Run("Error - Cannot reassign on merged PR", func(t *testing.T) {
		// Create team and PR
		teamName := fmt.Sprintf("merged-team-%d", time.Now().UnixNano())
		authorID := fmt.Sprintf("author-%d", time.Now().UnixNano())
		reviewerID := fmt.Sprintf("reviewer-%d", time.Now().UnixNano())

		createTeamBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": authorID, "username": "Author", "is_active": true},
				{"user_id": reviewerID, "username": "Reviewer", "is_active": true},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", createTeamBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		// Create PR
		prID := fmt.Sprintf("pr-%d", time.Now().UnixNano())
		createPRBody := map[string]interface{}{
			"pull_request_id":   prID,
			"pull_request_name": "Feature",
			"author_id":         authorID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/create", createPRBody)
		if err != nil {
			t.Fatalf("Failed to create PR: %v", err)
		}
		resp.Body.Close()

		// Merge the PR
		mergeBody := map[string]interface{}{
			"pull_request_id": prID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/merge", mergeBody)
		if err != nil {
			t.Fatalf("Failed to merge PR: %v", err)
		}
		resp.Body.Close()

		// Try to reassign after merge
		reassignBody := map[string]interface{}{
			"pull_request_id": prID,
			"old_user_id":     reviewerID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/reassign", reassignBody)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		assertStatusCode(t, resp, http.StatusConflict)
		assertErrorCode(t, resp, "PR_MERGED")
	})

	t.Run("Error - No candidate for replacement", func(t *testing.T) {
		// Create team with only 2 users (author + 1 reviewer)
		teamName := fmt.Sprintf("nocandidate-team-%d", time.Now().UnixNano())
		authorID := fmt.Sprintf("author-%d", time.Now().UnixNano())
		reviewerID := fmt.Sprintf("reviewer-%d", time.Now().UnixNano())

		createTeamBody := map[string]interface{}{
			"team_name": teamName,
			"members": []map[string]interface{}{
				{"user_id": authorID, "username": "Author", "is_active": true},
				{"user_id": reviewerID, "username": "Reviewer", "is_active": true},
			},
		}

		resp, err := doRequest(http.MethodPost, "/team/add", createTeamBody)
		if err != nil {
			t.Fatalf("Failed to create team: %v", err)
		}
		resp.Body.Close()

		// Create PR (will assign 1 reviewer)
		prID := fmt.Sprintf("pr-%d", time.Now().UnixNano())
		createPRBody := map[string]interface{}{
			"pull_request_id":   prID,
			"pull_request_name": "Feature",
			"author_id":         authorID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/create", createPRBody)
		if err != nil {
			t.Fatalf("Failed to create PR: %v", err)
		}
		resp.Body.Close()

		// Try to reassign the only reviewer (no candidates available)
		reassignBody := map[string]interface{}{
			"pull_request_id": prID,
			"old_user_id":     reviewerID,
		}

		resp, err = doRequest(http.MethodPost, "/pullRequest/reassign", reassignBody)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		assertStatusCode(t, resp, http.StatusConflict)
		assertErrorCode(t, resp, "NO_CANDIDATE")
	})

	t.Run("Error - PR not found", func(t *testing.T) {
		reassignBody := map[string]interface{}{
			"pull_request_id": "nonexistent-pr",
			"old_user_id":     "some-user",
		}

		resp, err := doRequest(http.MethodPost, "/pullRequest/reassign", reassignBody)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}

		assertStatusCode(t, resp, http.StatusNotFound)
		assertErrorCode(t, resp, "NOT_FOUND")
	})
}
