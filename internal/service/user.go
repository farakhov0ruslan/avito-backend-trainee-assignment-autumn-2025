package service

import (
	"context"
	"fmt"

	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models"
	"avito-backend-trainee-assignment-autumn-2025/internal/dto/request"
	"avito-backend-trainee-assignment-autumn-2025/internal/dto/response"
	"avito-backend-trainee-assignment-autumn-2025/internal/repository"
	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

// UserServiceImpl implements UserService
type UserServiceImpl struct {
	userRepo repository.UserRepository
	prRepo   repository.PRRepository
}

// NewUserService creates a new user service
func NewUserService(
	userRepo repository.UserRepository,
	prRepo repository.PRRepository,
) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

// SetUserActive sets the active status of a user
func (s *UserServiceImpl) SetUserActive(ctx context.Context, req *request.SetUserActiveRequest) (*response.SetUserActiveResponse, error) {
	// Validate input
	if req.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	logger.Info("Setting user %s active status to %t", req.UserID, req.IsActive)

	// Update user active status
	err := s.userRepo.SetActive(ctx, req.UserID, req.IsActive)
	if err != nil {
		logger.Error("Failed to set active status for user %s: %v", req.UserID, err)
		return nil, err
	}

	// Get updated user to return in response
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		logger.Error("Failed to get user %s after updating: %v", req.UserID, err)
		return nil, err
	}

	logger.Info("Successfully set user %s active status to %t", req.UserID, req.IsActive)

	// Convert to response DTO
	return &response.SetUserActiveResponse{
		User: convertUserToResponse(user),
	}, nil
}

// GetUserReviews retrieves all pull requests where the user is assigned as a reviewer
func (s *UserServiceImpl) GetUserReviews(ctx context.Context, userID string) (*response.GetUserReviewsResponse, error) {
	// Validate input
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	logger.Info("Retrieving reviews for user: %s", userID)

	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		logger.Error("Failed to get user %s: %v", userID, err)
		return nil, err
	}

	// Get all PRs where user is a reviewer
	prs, err := s.prRepo.GetPRsByReviewerID(ctx, userID)
	if err != nil {
		logger.Error("Failed to get PRs for reviewer %s: %v", userID, err)
		return nil, fmt.Errorf("failed to get PRs for reviewer: %w", err)
	}

	logger.Info("Successfully retrieved %d PRs for reviewer %s", len(prs), userID)

	// Convert to response DTO
	prResponses := make([]response.PullRequestShortResponse, 0, len(prs))
	for _, pr := range prs {
		prResponses = append(prResponses, response.PullRequestShortResponse{
			PullRequestID:   pr.ID,
			PullRequestName: pr.Name,
			AuthorID:        pr.AuthorID,
			Status:          string(pr.Status),
		})
	}

	return &response.GetUserReviewsResponse{
		UserID:       userID,
		PullRequests: prResponses,
	}, nil
}

// convertUserToResponse converts a User model to UserResponse DTO
func convertUserToResponse(user *models.User) response.UserResponse {
	return response.UserResponse{
		UserID:   user.ID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}
