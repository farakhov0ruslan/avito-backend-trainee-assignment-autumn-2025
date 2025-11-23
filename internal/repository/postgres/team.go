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

// TeamRepository implements repository.TeamRepository for PostgreSQL
type TeamRepository struct {
	pool *pgxpool.Pool
}

// NewTeamRepository creates a new team repository
func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pool: pool}
}

// Create creates a new team
func (r *TeamRepository) Create(ctx context.Context, team *models.Team) error {
	executor := repository.GetTx(ctx, r.pool)

	query := `INSERT INTO teams (name) VALUES ($1)`

	_, err := executor.Exec(ctx, query, team.Name)
	if err != nil {
		logger.Error("Failed to create team %s: %v", team.Name, err)
		// Check for unique violation
		if isPgUniqueViolation(err) {
			return pkgerrors.ErrTeamExists
		}
		return fmt.Errorf("failed to create team: %w", err)
	}

	logger.Info("Created team: %s", team.Name)
	return nil
}

// GetByName retrieves a team by name with all members
func (r *TeamRepository) GetByName(ctx context.Context, name string) (*models.Team, error) {
	executor := repository.GetTx(ctx, r.pool)

	// First, check if team exists
	var teamExists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1)`
	err := executor.QueryRow(ctx, checkQuery, name).Scan(&teamExists)
	if err != nil {
		logger.Error("Failed to check team existence %s: %v", name, err)
		return nil, fmt.Errorf("failed to check team existence: %w", err)
	}

	if !teamExists {
		return nil, pkgerrors.ErrTeamNotFound
	}

	// Get team members
	query := `
		SELECT id, username, team_name, is_active
		FROM users
		WHERE team_name = $1
		ORDER BY username
	`

	rows, err := executor.Query(ctx, query, name)
	if err != nil {
		logger.Error("Failed to get team members for %s: %v", name, err)
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}
	defer rows.Close()

	team := &models.Team{
		Name:    name,
		Members: make([]models.User, 0),
	}

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			logger.Error("Failed to scan user for team %s: %v", name, err)
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		team.Members = append(team.Members, user)
	}

	if err := rows.Err(); err != nil {
		logger.Error("Error iterating team members for %s: %v", name, err)
		return nil, fmt.Errorf("error iterating team members: %w", err)
	}

	logger.Debug("Retrieved team %s with %d members", name, len(team.Members))
	return team, nil
}

// Exists checks if a team exists
func (r *TeamRepository) Exists(ctx context.Context, name string) (bool, error) {
	executor := repository.GetTx(ctx, r.pool)

	query := `SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1)`

	var exists bool
	err := executor.QueryRow(ctx, query, name).Scan(&exists)
	if err != nil {
		logger.Error("Failed to check team existence %s: %v", name, err)
		return false, fmt.Errorf("failed to check team existence: %w", err)
	}

	return exists, nil
}
