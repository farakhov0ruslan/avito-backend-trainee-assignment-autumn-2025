package service

import (
	"context"
	"fmt"

	"avito-backend-trainee-assignment-autumn-2025/internal/domain/models"
	"avito-backend-trainee-assignment-autumn-2025/internal/dto/request"
	"avito-backend-trainee-assignment-autumn-2025/internal/dto/response"
	"avito-backend-trainee-assignment-autumn-2025/internal/repository"
	pkgerrors "avito-backend-trainee-assignment-autumn-2025/pkg/errors"
	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

// TeamServiceImpl implements TeamService
type TeamServiceImpl struct {
	teamRepo  repository.TeamRepository
	userRepo  repository.UserRepository
	txManager repository.TransactionManager
}

// NewTeamService creates a new team service
func NewTeamService(
	teamRepo repository.TeamRepository,
	userRepo repository.UserRepository,
	txManager repository.TransactionManager,
) *TeamServiceImpl {
	return &TeamServiceImpl{
		teamRepo:  teamRepo,
		userRepo:  userRepo,
		txManager: txManager,
	}
}

// CreateTeam creates a new team with members atomically
func (s *TeamServiceImpl) CreateTeam(ctx context.Context, req *request.CreateTeamRequest) (*response.CreateTeamResponse, error) {
	// Validate input
	if req.TeamName == "" {
		return nil, fmt.Errorf("team name is required")
	}

	logger.Info("Creating team: %s with %d members", req.TeamName, len(req.Members))

	// Check if team already exists
	exists, err := s.teamRepo.Exists(ctx, req.TeamName)
	if err != nil {
		logger.Error("Failed to check team existence: %v", err)
		return nil, fmt.Errorf("failed to check team existence: %w", err)
	}
	if exists {
		logger.Warn("Team already exists: %s", req.TeamName)
		return nil, pkgerrors.ErrTeamExists
	}

	// Create team and members atomically in a transaction
	var createdTeam *models.Team
	err = s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// Create team
		team := &models.Team{
			Name:    req.TeamName,
			Members: []models.User{},
		}

		if err := s.teamRepo.Create(txCtx, team); err != nil {
			logger.Error("Failed to create team %s: %v", req.TeamName, err)
			return fmt.Errorf("failed to create team: %w", err)
		}

		// Create all members
		for _, memberReq := range req.Members {
			// Validate member
			if memberReq.UserID == "" {
				return fmt.Errorf("user_id is required for all members")
			}
			if memberReq.Username == "" {
				return fmt.Errorf("username is required for all members")
			}

			user := &models.User{
				ID:       memberReq.UserID,
				Username: memberReq.Username,
				TeamName: req.TeamName,
				IsActive: memberReq.IsActive,
			}

			if err := s.userRepo.Create(txCtx, user); err != nil {
				logger.Error("Failed to create user %s for team %s: %v", user.ID, req.TeamName, err)
				return fmt.Errorf("failed to create user %s: %w", user.ID, err)
			}

			team.Members = append(team.Members, *user)
		}

		createdTeam = team
		return nil
	})

	if err != nil {
		return nil, err
	}

	logger.Info("Successfully created team %s with %d members", req.TeamName, len(createdTeam.Members))

	// Convert to response DTO
	return &response.CreateTeamResponse{
		Team: convertTeamToResponse(createdTeam),
	}, nil
}

// GetTeam retrieves a team with all its members
func (s *TeamServiceImpl) GetTeam(ctx context.Context, teamName string) (*response.TeamResponse, error) {
	// Validate input
	if teamName == "" {
		return nil, fmt.Errorf("team name is required")
	}

	logger.Info("Retrieving team: %s", teamName)

	// Get team from repository
	team, err := s.teamRepo.GetByName(ctx, teamName)
	if err != nil {
		logger.Error("Failed to get team %s: %v", teamName, err)
		return nil, err
	}

	logger.Info("Successfully retrieved team %s with %d members", teamName, len(team.Members))

	// Convert to response DTO
	return convertTeamToResponsePtr(team), nil
}

// convertTeamToResponse converts a Team model to TeamResponse DTO
func convertTeamToResponse(team *models.Team) response.TeamResponse {
	members := make([]response.TeamMemberResponse, 0, len(team.Members))
	for _, member := range team.Members {
		members = append(members, response.TeamMemberResponse{
			UserID:   member.ID,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}

	return response.TeamResponse{
		TeamName: team.Name,
		Members:  members,
	}
}

// convertTeamToResponsePtr is a helper that returns a pointer to TeamResponse
func convertTeamToResponsePtr(team *models.Team) *response.TeamResponse {
	resp := convertTeamToResponse(team)
	return &resp
}
