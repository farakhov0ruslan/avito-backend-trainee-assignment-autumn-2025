package service

import (
	"context"

	"avito-backend-trainee-assignment-autumn-2025/internal/dto/request"
	"avito-backend-trainee-assignment-autumn-2025/internal/dto/response"
)

// TeamService defines business logic for team operations
type TeamService interface {
	// CreateTeam creates a new team with members
	// Returns error if team already exists or validation fails
	CreateTeam(ctx context.Context, req *request.CreateTeamRequest) (*response.CreateTeamResponse, error)

	// GetTeam retrieves a team with all its members
	// Returns error if team doesn't exist
	GetTeam(ctx context.Context, teamName string) (*response.TeamResponse, error)
}

// UserService defines business logic for user operations
type UserService interface {
	// SetUserActive sets the active status of a user
	// Returns error if user doesn't exist
	SetUserActive(ctx context.Context, req *request.SetUserActiveRequest) (*response.SetUserActiveResponse, error)

	// GetUserReviews retrieves all pull requests where the user is assigned as a reviewer
	// Returns error if user doesn't exist
	GetUserReviews(ctx context.Context, userID string) (*response.GetUserReviewsResponse, error)
}

// PRService defines business logic for pull request operations
type PRService interface {
	// CreatePR creates a new pull request and automatically assigns up to 2 reviewers
	// Reviewers are selected randomly from the author's team (excluding the author)
	// Only active users can be assigned as reviewers
	// Returns error if PR already exists, author doesn't exist, or validation fails
	CreatePR(ctx context.Context, req *request.CreatePRRequest) (*response.CreatePRResponse, error)

	// MergePR merges a pull request (sets status to MERGED)
	// This operation is idempotent - if already merged, returns current state
	// Returns error if PR doesn't exist
	MergePR(ctx context.Context, req *request.MergePRRequest) (*response.MergePRResponse, error)

	// ReassignReviewer replaces one reviewer with another random active member
	// The new reviewer is selected from the replaced reviewer's team
	// Returns error if:
	// - PR doesn't exist
	// - PR is already merged (PR_MERGED)
	// - old_user_id is not assigned as reviewer (NOT_ASSIGNED)
	// - No suitable candidates available (NO_CANDIDATE)
	ReassignReviewer(ctx context.Context, req *request.ReassignReviewerRequest) (*response.ReassignReviewerResponse, error)
}
