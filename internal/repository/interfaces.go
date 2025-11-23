package repository

import (
	"context"

	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models"
)

// TeamRepository defines methods for working with teams
type TeamRepository interface {
	Create(ctx context.Context, team *models.Team) error
	GetByName(ctx context.Context, name string) (*models.Team, error)
	Exists(ctx context.Context, name string) (bool, error)
}

// UserRepository defines methods for working with users
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByTeamName(ctx context.Context, teamName string) ([]models.User, error)
	SetActive(ctx context.Context, userID string, isActive bool) error
}

// PRRepository defines methods for working with pull requests
type PRRepository interface {
	Create(ctx context.Context, pr *models.PullRequest) error
	GetByID(ctx context.Context, id string) (*models.PullRequest, error)
	Update(ctx context.Context, pr *models.PullRequest) error
	Merge(ctx context.Context, prID string) (*models.PullRequest, error)
	GetReviewersByPRID(ctx context.Context, prID string) ([]string, error)
	AddReviewer(ctx context.Context, prID, reviewerID string) error
	RemoveReviewer(ctx context.Context, prID, reviewerID string) error
	GetPRsByReviewerID(ctx context.Context, reviewerID string) ([]models.PullRequest, error)
}
