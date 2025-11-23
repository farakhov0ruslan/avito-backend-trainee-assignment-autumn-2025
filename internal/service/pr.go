package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models"
	"avito-backend-trainee-assignment-autumn-2025/internal/dto/request"
	"avito-backend-trainee-assignment-autumn-2025/internal/dto/response"
	"avito-backend-trainee-assignment-autumn-2025/internal/repository"
	pkgerrors "avito-backend-trainee-assignment-autumn-2025/pkg/errors"
	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

// PRServiceImpl implements PRService
type PRServiceImpl struct {
	prRepo    repository.PRRepository
	userRepo  repository.UserRepository
	teamRepo  repository.TeamRepository
	txManager repository.TransactionManager
	rand      *rand.Rand
}

// NewPRService creates a new PR service
func NewPRService(
	prRepo repository.PRRepository,
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
	txManager repository.TransactionManager,
) *PRServiceImpl {
	// Initialize random number generator with current time as seed
	source := rand.NewSource(time.Now().UnixNano())
	return &PRServiceImpl{
		prRepo:    prRepo,
		userRepo:  userRepo,
		teamRepo:  teamRepo,
		txManager: txManager,
		rand:      rand.New(source),
	}
}

// CreatePR creates a new pull request and automatically assigns up to 2 reviewers
func (s *PRServiceImpl) CreatePR(ctx context.Context, req *request.CreatePRRequest) (*response.CreatePRResponse, error) {
	// Validate input
	if req.PullRequestID == "" {
		return nil, fmt.Errorf("pull_request_id is required")
	}
	if req.PullRequestName == "" {
		return nil, fmt.Errorf("pull_request_name is required")
	}
	if req.AuthorID == "" {
		return nil, fmt.Errorf("author_id is required")
	}

	logger.Info("Creating PR: %s (author: %s)", req.PullRequestID, req.AuthorID)

	// Check if author exists and get their team
	author, err := s.userRepo.GetByID(ctx, req.AuthorID)
	if err != nil {
		logger.Error("Failed to get author %s: %v", req.AuthorID, err)
		return nil, err
	}

	// Get author's team
	team, err := s.teamRepo.GetByName(ctx, author.TeamName)
	if err != nil {
		logger.Error("Failed to get team %s: %v", author.TeamName, err)
		return nil, fmt.Errorf("failed to get author's team: %w", err)
	}

	// Get active members excluding the author
	candidates := team.GetActiveMembersExcept(req.AuthorID)
	logger.Debug("Found %d active candidates for PR %s", len(candidates), req.PullRequestID)

	// Select up to 2 random reviewers
	selectedReviewers := s.selectRandomReviewers(candidates, 2)
	reviewerIDs := make([]string, len(selectedReviewers))
	for i, reviewer := range selectedReviewers {
		reviewerIDs[i] = reviewer.ID
	}

	logger.Info("Selected %d reviewers for PR %s: %v", len(reviewerIDs), req.PullRequestID, reviewerIDs)

	// Create PR and assign reviewers in a transaction
	var createdPR *models.PullRequest
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Create PR
		pr := &models.PullRequest{
			ID:                req.PullRequestID,
			Name:              req.PullRequestName,
			AuthorID:          req.AuthorID,
			Status:            models.PRStatusOpen,
			AssignedReviewers: []string{},
			CreatedAt:         time.Now(),
		}

		if err := s.prRepo.Create(txCtx, pr); err != nil {
			logger.Error("Failed to create PR %s: %v", req.PullRequestID, err)
			return err
		}

		// Assign reviewers
		for _, reviewerID := range reviewerIDs {
			if err := s.prRepo.AddReviewer(txCtx, req.PullRequestID, reviewerID); err != nil {
				logger.Error("Failed to assign reviewer %s to PR %s: %v", reviewerID, req.PullRequestID, err)
				return fmt.Errorf("failed to assign reviewer %s: %w", reviewerID, err)
			}
		}

		// Get the created PR with reviewers
		pr.AssignedReviewers = reviewerIDs
		createdPR = pr
		return nil
	})

	if err != nil {
		return nil, err
	}

	logger.Info("Successfully created PR %s with %d reviewers", req.PullRequestID, len(reviewerIDs))

	// Convert to response DTO
	return &response.CreatePRResponse{
		PR: convertPRToResponse(createdPR),
	}, nil
}

// MergePR merges a pull request (idempotent operation)
func (s *PRServiceImpl) MergePR(ctx context.Context, req *request.MergePRRequest) (*response.MergePRResponse, error) {
	// Validate input
	if req.PullRequestID == "" {
		return nil, fmt.Errorf("pull_request_id is required")
	}

	logger.Info("Merging PR: %s", req.PullRequestID)

	// Get PR to check if it exists
	pr, err := s.prRepo.GetByID(ctx, req.PullRequestID)
	if err != nil {
		logger.Error("Failed to get PR %s: %v", req.PullRequestID, err)
		return nil, err
	}

	// Check if already merged (idempotency)
	if pr.IsMerged() {
		logger.Info("PR %s is already merged, returning current state", req.PullRequestID)
		return &response.MergePRResponse{
			PR: convertPRToResponse(pr),
		}, nil
	}

	// Merge PR
	mergedPR, err := s.prRepo.Merge(ctx, req.PullRequestID)
	if err != nil {
		logger.Error("Failed to merge PR %s: %v", req.PullRequestID, err)
		return nil, err
	}

	logger.Info("Successfully merged PR %s", req.PullRequestID)

	// Convert to response DTO
	return &response.MergePRResponse{
		PR: convertPRToResponse(mergedPR),
	}, nil
}

// ReassignReviewer replaces one reviewer with another random active member
func (s *PRServiceImpl) ReassignReviewer(ctx context.Context, req *request.ReassignReviewerRequest) (*response.ReassignReviewerResponse, error) {
	// Validate input
	if req.PullRequestID == "" {
		return nil, fmt.Errorf("pull_request_id is required")
	}
	if req.OldUserID == "" {
		return nil, fmt.Errorf("old_user_id is required")
	}

	logger.Info("Reassigning reviewer for PR %s: replacing %s", req.PullRequestID, req.OldUserID)

	// Get PR
	pr, err := s.prRepo.GetByID(ctx, req.PullRequestID)
	if err != nil {
		logger.Error("Failed to get PR %s: %v", req.PullRequestID, err)
		return nil, err
	}

	// Check if PR is merged
	if pr.IsMerged() {
		logger.Warn("Cannot reassign reviewer for merged PR %s", req.PullRequestID)
		return nil, pkgerrors.ErrPRMerged
	}

	// Check if old_user_id is assigned as reviewer
	if !pr.IsReviewerAssigned(req.OldUserID) {
		logger.Warn("User %s is not assigned as reviewer for PR %s", req.OldUserID, req.PullRequestID)
		return nil, pkgerrors.ErrReviewerNotAssigned
	}

	// Get the team of the user being replaced
	oldUser, err := s.userRepo.GetByID(ctx, req.OldUserID)
	if err != nil {
		logger.Error("Failed to get user %s: %v", req.OldUserID, err)
		return nil, err
	}

	team, err := s.teamRepo.GetByName(ctx, oldUser.TeamName)
	if err != nil {
		logger.Error("Failed to get team %s: %v", oldUser.TeamName, err)
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	// Get active candidates excluding the old reviewer, the PR author, and other current reviewers
	candidates := []models.User{}
	for _, member := range team.GetActiveMembers() {
		// Skip if member is the old reviewer or the PR author
		if member.ID == req.OldUserID || member.ID == pr.AuthorID {
			continue
		}

		// Skip if member is already assigned as a reviewer (excluding the one being replaced)
		isCurrentlyAssigned := false
		for _, assignedReviewerID := range pr.AssignedReviewers {
			if member.ID == assignedReviewerID {
				isCurrentlyAssigned = true
				break
			}
		}

		if !isCurrentlyAssigned {
			candidates = append(candidates, member)
		}
	}

	logger.Debug("Found %d candidates for reassignment in team %s", len(candidates), team.Name)

	// Check if there are any candidates
	if len(candidates) == 0 {
		logger.Warn("No candidates available for reassignment in team %s", team.Name)
		return nil, pkgerrors.ErrNoCandidates
	}

	// Select random candidate
	selectedReviewers := s.selectRandomReviewers(candidates, 1)
	if len(selectedReviewers) == 0 {
		logger.Warn("Failed to select reviewer from %d candidates", len(candidates))
		return nil, pkgerrors.ErrNoCandidates
	}

	newReviewerID := selectedReviewers[0].ID
	logger.Info("Selected new reviewer %s to replace %s for PR %s", newReviewerID, req.OldUserID, req.PullRequestID)

	// Replace reviewer in a transaction
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Remove old reviewer
		if err := s.prRepo.RemoveReviewer(txCtx, req.PullRequestID, req.OldUserID); err != nil {
			logger.Error("Failed to remove reviewer %s from PR %s: %v", req.OldUserID, req.PullRequestID, err)
			return fmt.Errorf("failed to remove old reviewer: %w", err)
		}

		// Add new reviewer
		if err := s.prRepo.AddReviewer(txCtx, req.PullRequestID, newReviewerID); err != nil {
			logger.Error("Failed to add reviewer %s to PR %s: %v", newReviewerID, req.PullRequestID, err)
			return fmt.Errorf("failed to add new reviewer: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Get updated PR
	updatedPR, err := s.prRepo.GetByID(ctx, req.PullRequestID)
	if err != nil {
		logger.Error("Failed to get updated PR %s: %v", req.PullRequestID, err)
		return nil, fmt.Errorf("failed to get updated PR: %w", err)
	}

	logger.Info("Successfully reassigned reviewer for PR %s: %s -> %s", req.PullRequestID, req.OldUserID, newReviewerID)

	// Convert to response DTO
	return &response.ReassignReviewerResponse{
		PR:         convertPRToResponse(updatedPR),
		ReplacedBy: newReviewerID,
	}, nil
}

// selectRandomReviewers selects up to maxCount random reviewers from candidates
func (s *PRServiceImpl) selectRandomReviewers(candidates []models.User, maxCount int) []models.User {
	if len(candidates) == 0 {
		return []models.User{}
	}

	// Determine how many to select
	count := maxCount
	if len(candidates) < count {
		count = len(candidates)
	}

	// Create a copy of candidates to avoid modifying original slice
	available := make([]models.User, len(candidates))
	copy(available, candidates)

	// Shuffle and select first 'count' elements (Fisher-Yates shuffle)
	selected := make([]models.User, count)
	for i := 0; i < count; i++ {
		// Pick random index from remaining elements
		j := i + s.rand.Intn(len(available)-i)
		// Swap
		available[i], available[j] = available[j], available[i]
		// Add to selected
		selected[i] = available[i]
	}

	return selected
}

// convertPRToResponse converts a PullRequest model to PullRequestResponse DTO
func convertPRToResponse(pr *models.PullRequest) response.PullRequestResponse {
	return response.PullRequestResponse{
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         &pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}
