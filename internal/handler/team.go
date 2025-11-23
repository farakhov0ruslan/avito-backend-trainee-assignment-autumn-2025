package handler

import (
	"fmt"
	"net/http"

	"avito-backend-trainee-assignment-autumn-2025/internal/dto/request"
	"avito-backend-trainee-assignment-autumn-2025/internal/service"
	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

// TeamHandler handles team-related HTTP requests
type TeamHandler struct {
	teamService service.TeamService
}

// NewTeamHandler creates a new team handler
func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

// CreateTeam handles POST /team/add
func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req request.CreateTeamRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		respondWithError(w, fmt.Errorf("invalid request body: %w", err))
		return
	}

	// Validate input
	if req.TeamName == "" {
		respondWithError(w, fmt.Errorf("team_name is required"))
		return
	}
	if len(req.Members) == 0 {
		respondWithError(w, fmt.Errorf("members array cannot be empty"))
		return
	}

	logger.Info("Creating team: %s", req.TeamName)

	// Call service
	resp, err := h.teamService.CreateTeam(r.Context(), &req)
	if err != nil {
		logger.Error("Failed to create team %s: %v", req.TeamName, err)
		respondWithError(w, err)
		return
	}

	// Send response
	respondWithJSON(w, http.StatusCreated, resp)
}

// GetTeam handles GET /team/get?team_name=...
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	// Get team_name from query parameters
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		respondWithError(w, fmt.Errorf("team_name query parameter is required"))
		return
	}

	logger.Info("Getting team: %s", teamName)

	// Call service
	resp, err := h.teamService.GetTeam(r.Context(), teamName)
	if err != nil {
		logger.Error("Failed to get team %s: %v", teamName, err)
		respondWithError(w, err)
		return
	}

	// Send response
	respondWithJSON(w, http.StatusOK, resp)
}
