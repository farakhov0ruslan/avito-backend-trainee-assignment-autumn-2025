package handler

import (
	"fmt"
	"net/http"

	"avito-backend-trainee-assignment-autumn-2025/internal/dto/request"
	"avito-backend-trainee-assignment-autumn-2025/internal/service"
	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// SetIsActive handles POST /users/setIsActive
func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req request.SetUserActiveRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		respondWithError(w, fmt.Errorf("invalid request body: %w", err))
		return
	}

	// Validate input
	if req.UserID == "" {
		respondWithError(w, fmt.Errorf("user_id is required"))
		return
	}

	logger.Info("Setting user %s active status to %t", req.UserID, req.IsActive)

	// Call service
	resp, err := h.userService.SetUserActive(r.Context(), &req)
	if err != nil {
		logger.Error("Failed to set active status for user %s: %v", req.UserID, err)
		respondWithError(w, err)
		return
	}

	// Send response
	respondWithJSON(w, http.StatusOK, resp)
}

// GetUserReviews handles GET /users/getReview?user_id=...
func (h *UserHandler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	// Get user_id from query parameters
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		respondWithError(w, fmt.Errorf("user_id query parameter is required"))
		return
	}

	logger.Info("Getting reviews for user: %s", userID)

	// Call service
	resp, err := h.userService.GetUserReviews(r.Context(), userID)
	if err != nil {
		logger.Error("Failed to get reviews for user %s: %v", userID, err)
		respondWithError(w, err)
		return
	}

	// Send response
	respondWithJSON(w, http.StatusOK, resp)
}
