package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models"
	"avito-backend-trainee-assignment-autumn-2025/internal/repository"
	pkgerrors "avito-backend-trainee-assignment-autumn-2025/pkg/errors"
	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

// UserRepository implements repository.UserRepository for PostgreSQL
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository creates a new user repository
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	executor := repository.GetTx(ctx, r.pool)

	query := `
		INSERT INTO users (id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
	`

	_, err := executor.Exec(ctx, query, user.ID, user.Username, user.TeamName, user.IsActive)
	if err != nil {
		logger.Error("Failed to create user %s: %v", user.ID, err)
		// Check for unique violation
		if isPgUniqueViolation(err) {
			return pkgerrors.ErrUserAlreadyExists
		}
		// Check for foreign key violation (team doesn't exist)
		if isPgForeignKeyViolation(err) {
			return pkgerrors.ErrTeamNotFound
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	logger.Info("Created user: %s (username: %s, team: %s)", user.ID, user.Username, user.TeamName)
	return nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	executor := repository.GetTx(ctx, r.pool)

	query := `
		UPDATE users
		SET username = $2, team_name = $3, is_active = $4, updated_at = NOW()
		WHERE id = $1
	`

	commandTag, err := executor.Exec(ctx, query, user.ID, user.Username, user.TeamName, user.IsActive)
	if err != nil {
		logger.Error("Failed to update user %s: %v", user.ID, err)
		// Check for foreign key violation (team doesn't exist)
		if isPgForeignKeyViolation(err) {
			return pkgerrors.ErrTeamNotFound
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return pkgerrors.ErrUserNotFound
	}

	logger.Info("Updated user: %s", user.ID)
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	executor := repository.GetTx(ctx, r.pool)

	query := `
		SELECT id, username, team_name, is_active
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := executor.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		if isPgNoRows(err) {
			return nil, pkgerrors.ErrUserNotFound
		}
		logger.Error("Failed to get user %s: %v", id, err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	logger.Debug("Retrieved user: %s", user.ID)
	return &user, nil
}

// GetByTeamName retrieves all users in a team
func (r *UserRepository) GetByTeamName(ctx context.Context, teamName string) ([]models.User, error) {
	executor := repository.GetTx(ctx, r.pool)

	query := `
		SELECT id, username, team_name, is_active
		FROM users
		WHERE team_name = $1
		ORDER BY username
	`

	rows, err := executor.Query(ctx, query, teamName)
	if err != nil {
		logger.Error("Failed to get users for team %s: %v", teamName, err)
		return nil, fmt.Errorf("failed to get users by team: %w", err)
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			logger.Error("Failed to scan user for team %s: %v", teamName, err)
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Error iterating users for team %s: %v", teamName, err)
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	logger.Debug("Retrieved %d users for team %s", len(users), teamName)
	return users, nil
}

// SetActive sets the active status of a user
func (r *UserRepository) SetActive(ctx context.Context, userID string, isActive bool) error {
	executor := repository.GetTx(ctx, r.pool)

	query := `
		UPDATE users
		SET is_active = $2, updated_at = NOW()
		WHERE id = $1
	`

	commandTag, err := executor.Exec(ctx, query, userID, isActive)
	if err != nil {
		logger.Error("Failed to set active status for user %s: %v", userID, err)
		return fmt.Errorf("failed to set active status: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return pkgerrors.ErrUserNotFound
	}

	logger.Info("Set user %s active status to %t", userID, isActive)
	return nil
}
