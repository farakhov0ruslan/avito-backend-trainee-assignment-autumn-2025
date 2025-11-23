package handler

import (
	"fmt"
	"net/http"

	"avito-backend-trainee-assignment-autumn-2025/internal/dto/request"
	"avito-backend-trainee-assignment-autumn-2025/internal/service"
	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

// PRHandler handles pull request-related HTTP requests
type PRHandler struct {
	prService service.PRService
}

// NewPRHandler creates a new PR handler
func NewPRHandler(prService service.PRService) *PRHandler {
	return &PRHandler{
		prService: prService,
	}
}

// CreatePR handles POST /pullRequest/create
func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req request.CreatePRRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		respondWithError(w, fmt.Errorf("invalid request body: %w", err))
		return
	}

	// Validate input
	if req.PullRequestID == "" {
		respondWithError(w, fmt.Errorf("pull_request_id is required"))
		return
	}
	if req.PullRequestName == "" {
		respondWithError(w, fmt.Errorf("pull_request_name is required"))
		return
	}
	if req.AuthorID == "" {
		respondWithError(w, fmt.Errorf("author_id is required"))
		return
	}

	logger.Info("Creating PR: %s (author: %s)", req.PullRequestID, req.AuthorID)

	// Call service
	resp, err := h.prService.CreatePR(r.Context(), &req)
	if err != nil {
		logger.Error("Failed to create PR %s: %v", req.PullRequestID, err)
		respondWithError(w, err)
		return
	}

	// Send response
	respondWithJSON(w, http.StatusCreated, resp)
}

// MergePR handles POST /pullRequest/merge
func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req request.MergePRRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		respondWithError(w, fmt.Errorf("invalid request body: %w", err))
		return
	}

	// Validate input
	if req.PullRequestID == "" {
		respondWithError(w, fmt.Errorf("pull_request_id is required"))
		return
	}

	logger.Info("Merging PR: %s", req.PullRequestID)

	// Call service
	resp, err := h.prService.MergePR(r.Context(), &req)
	if err != nil {
		logger.Error("Failed to merge PR %s: %v", req.PullRequestID, err)
		respondWithError(w, err)
		return
	}

	// Send response (200 OK for idempotent operation)
	respondWithJSON(w, http.StatusOK, resp)
}

// ReassignReviewer handles POST /pullRequest/reassign
func (h *PRHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req request.ReassignReviewerRequest
	if err := decodeJSONBody(w, r, &req); err != nil {
		respondWithError(w, fmt.Errorf("invalid request body: %w", err))
		return
	}

	// Validate input
	if req.PullRequestID == "" {
		respondWithError(w, fmt.Errorf("pull_request_id is required"))
		return
	}
	if req.OldUserID == "" {
		respondWithError(w, fmt.Errorf("old_user_id is required"))
		return
	}

	logger.Info("Reassigning reviewer for PR %s: replacing %s", req.PullRequestID, req.OldUserID)

	// Call service
	resp, err := h.prService.ReassignReviewer(r.Context(), &req)
	if err != nil {
		logger.Error("Failed to reassign reviewer for PR %s: %v", req.PullRequestID, err)
		respondWithError(w, err)
		return
	}

	// Send response
	respondWithJSON(w, http.StatusOK, resp)
}
