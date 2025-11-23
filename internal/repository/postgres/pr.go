package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models"
	"avito-backend-trainee-assignment-autumn-2025/internal/repository"
	pkgerrors "avito-backend-trainee-assignment-autumn-2025/pkg/errors"
	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

// PRRepository implements repository.PRRepository for PostgreSQL
type PRRepository struct {
	pool *pgxpool.Pool
}

// NewPRRepository creates a new pull request repository
func NewPRRepository(pool *pgxpool.Pool) *PRRepository {
	return &PRRepository{pool: pool}
}

// Create creates a new pull request
func (r *PRRepository) Create(ctx context.Context, pr *models.PullRequest) error {
	executor := repository.GetTx(ctx, r.pool)

	query := `
		INSERT INTO pull_requests (id, name, author_id, status)
		VALUES ($1, $2, $3, $4)
	`

	_, err := executor.Exec(ctx, query, pr.ID, pr.Name, pr.AuthorID, pr.Status)
	if err != nil {
		logger.Error("Failed to create PR %s: %v", pr.ID, err)
		// Check for unique violation
		if isPgUniqueViolation(err) {
			return pkgerrors.ErrPRExists
		}
		// Check for foreign key violation (author doesn't exist)
		if isPgForeignKeyViolation(err) {
			return pkgerrors.ErrUserNotFound
		}
		// Check for check constraint violation (invalid status)
		if isPgCheckViolation(err) {
			return fmt.Errorf("invalid PR status: %s", pr.Status)
		}
		return fmt.Errorf("failed to create PR: %w", err)
	}

	logger.Info("Created PR: %s (name: %s, author: %s)", pr.ID, pr.Name, pr.AuthorID)
	return nil
}

// GetByID retrieves a pull request by ID with all reviewers
func (r *PRRepository) GetByID(ctx context.Context, id string) (*models.PullRequest, error) {
	executor := repository.GetTx(ctx, r.pool)

	// Get PR details
	query := `
		SELECT id, name, author_id, status, created_at, merged_at
		FROM pull_requests
		WHERE id = $1
	`

	var pr models.PullRequest
	err := executor.QueryRow(ctx, query, id).Scan(
		&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt,
	)
	if err != nil {
		if isPgNoRows(err) {
			return nil, pkgerrors.ErrPRNotFound
		}
		logger.Error("Failed to get PR %s: %v", id, err)
		return nil, fmt.Errorf("failed to get PR: %w", err)
	}

	// Get reviewers
	reviewers, err := r.GetReviewersByPRID(ctx, id)
	if err != nil {
		logger.Error("Failed to get reviewers for PR %s: %v", id, err)
		return nil, fmt.Errorf("failed to get reviewers: %w", err)
	}

	pr.AssignedReviewers = reviewers

	logger.Debug("Retrieved PR %s with %d reviewers", pr.ID, len(pr.AssignedReviewers))
	return &pr, nil
}

// Update updates an existing pull request
func (r *PRRepository) Update(ctx context.Context, pr *models.PullRequest) error {
	executor := repository.GetTx(ctx, r.pool)

	query := `
		UPDATE pull_requests
		SET name = $2, author_id = $3, status = $4, updated_at = NOW()
		WHERE id = $1
	`

	commandTag, err := executor.Exec(ctx, query, pr.ID, pr.Name, pr.AuthorID, pr.Status)
	if err != nil {
		logger.Error("Failed to update PR %s: %v", pr.ID, err)
		// Check for foreign key violation (author doesn't exist)
		if isPgForeignKeyViolation(err) {
			return pkgerrors.ErrUserNotFound
		}
		// Check for check constraint violation (invalid status)
		if isPgCheckViolation(err) {
			return fmt.Errorf("invalid PR status: %s", pr.Status)
		}
		return fmt.Errorf("failed to update PR: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return pkgerrors.ErrPRNotFound
	}

	logger.Info("Updated PR: %s", pr.ID)
	return nil
}

// Merge merges a pull request (sets status to MERGED and merged_at timestamp)
func (r *PRRepository) Merge(ctx context.Context, prID string) (*models.PullRequest, error) {
	executor := repository.GetTx(ctx, r.pool)

	// First check if PR exists and is not already merged
	var currentStatus models.PRStatus
	checkQuery := `SELECT status FROM pull_requests WHERE id = $1`
	err := executor.QueryRow(ctx, checkQuery, prID).Scan(&currentStatus)
	if err != nil {
		if isPgNoRows(err) {
			return nil, pkgerrors.ErrPRNotFound
		}
		logger.Error("Failed to check PR status %s: %v", prID, err)
		return nil, fmt.Errorf("failed to check PR status: %w", err)
	}

	if currentStatus == models.PRStatusMerged {
		return nil, pkgerrors.ErrPRMerged
	}

	// Update PR to merged status
	query := `
		UPDATE pull_requests
		SET status = $2, merged_at = $3, updated_at = NOW()
		WHERE id = $1
	`

	mergedAt := time.Now()
	_, err = executor.Exec(ctx, query, prID, models.PRStatusMerged, mergedAt)
	if err != nil {
		logger.Error("Failed to merge PR %s: %v", prID, err)
		return nil, fmt.Errorf("failed to merge PR: %w", err)
	}

	logger.Info("Merged PR: %s", prID)

	// Return updated PR
	return r.GetByID(ctx, prID)
}

// GetReviewersByPRID retrieves all reviewer IDs for a pull request
func (r *PRRepository) GetReviewersByPRID(ctx context.Context, prID string) ([]string, error) {
	executor := repository.GetTx(ctx, r.pool)

	query := `
		SELECT reviewer_id
		FROM pr_reviewers
		WHERE pr_id = $1
		ORDER BY assigned_at
	`

	rows, err := executor.Query(ctx, query, prID)
	if err != nil {
		logger.Error("Failed to get reviewers for PR %s: %v", prID, err)
		return nil, fmt.Errorf("failed to get reviewers: %w", err)
	}
	defer rows.Close()

	reviewers := make([]string, 0)
	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			logger.Error("Failed to scan reviewer for PR %s: %v", prID, err)
			return nil, fmt.Errorf("failed to scan reviewer: %w", err)
		}
		reviewers = append(reviewers, reviewerID)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Error iterating reviewers for PR %s: %v", prID, err)
		return nil, fmt.Errorf("error iterating reviewers: %w", err)
	}

	logger.Debug("Retrieved %d reviewers for PR %s", len(reviewers), prID)
	return reviewers, nil
}

// AddReviewer adds a reviewer to a pull request
func (r *PRRepository) AddReviewer(ctx context.Context, prID, reviewerID string) error {
	executor := repository.GetTx(ctx, r.pool)

	query := `
		INSERT INTO pr_reviewers (pr_id, reviewer_id)
		VALUES ($1, $2)
	`

	_, err := executor.Exec(ctx, query, prID, reviewerID)
	if err != nil {
		logger.Error("Failed to add reviewer %s to PR %s: %v", reviewerID, prID, err)
		// Check for unique violation (reviewer already assigned)
		if isPgUniqueViolation(err) {
			return fmt.Errorf("reviewer already assigned to this PR")
		}
		// Check for foreign key violation (PR or reviewer doesn't exist)
		if isPgForeignKeyViolation(err) {
			return pkgerrors.ErrPRNotFound
		}
		return fmt.Errorf("failed to add reviewer: %w", err)
	}

	logger.Info("Added reviewer %s to PR %s", reviewerID, prID)
	return nil
}

// RemoveReviewer removes a reviewer from a pull request
func (r *PRRepository) RemoveReviewer(ctx context.Context, prID, reviewerID string) error {
	executor := repository.GetTx(ctx, r.pool)

	query := `
		DELETE FROM pr_reviewers
		WHERE pr_id = $1 AND reviewer_id = $2
	`

	commandTag, err := executor.Exec(ctx, query, prID, reviewerID)
	if err != nil {
		logger.Error("Failed to remove reviewer %s from PR %s: %v", reviewerID, prID, err)
		return fmt.Errorf("failed to remove reviewer: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return pkgerrors.ErrReviewerNotAssigned
	}

	logger.Info("Removed reviewer %s from PR %s", reviewerID, prID)
	return nil
}

// GetPRsByReviewerID retrieves all pull requests assigned to a reviewer
func (r *PRRepository) GetPRsByReviewerID(ctx context.Context, reviewerID string) ([]models.PullRequest, error) {
	executor := repository.GetTx(ctx, r.pool)

	query := `
		SELECT DISTINCT pr.id, pr.name, pr.author_id, pr.status, pr.created_at, pr.merged_at
		FROM pull_requests pr
		INNER JOIN pr_reviewers prr ON pr.id = prr.pr_id
		WHERE prr.reviewer_id = $1
		ORDER BY pr.created_at DESC
	`

	rows, err := executor.Query(ctx, query, reviewerID)
	if err != nil {
		logger.Error("Failed to get PRs for reviewer %s: %v", reviewerID, err)
		return nil, fmt.Errorf("failed to get PRs by reviewer: %w", err)
	}
	defer rows.Close()

	prs := make([]models.PullRequest, 0)
	for rows.Next() {
		var pr models.PullRequest
		if err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt); err != nil {
			logger.Error("Failed to scan PR for reviewer %s: %v", reviewerID, err)
			return nil, fmt.Errorf("failed to scan PR: %w", err)
		}

		// Get reviewers for this PR
		reviewers, err := r.GetReviewersByPRID(ctx, pr.ID)
		if err != nil {
			logger.Error("Failed to get reviewers for PR %s: %v", pr.ID, err)
			return nil, fmt.Errorf("failed to get reviewers: %w", err)
		}
		pr.AssignedReviewers = reviewers

		prs = append(prs, pr)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Error iterating PRs for reviewer %s: %v", reviewerID, err)
		return nil, fmt.Errorf("error iterating PRs: %w", err)
	}

	logger.Debug("Retrieved %d PRs for reviewer %s", len(prs), reviewerID)
	return prs, nil
}
