package errors

import (
	"avito-backend-trainee-assignment-autumn-2025/internal/dto/response"
	"errors"
	"net/http"
)

// Domain errors
var (
	// Team errors
	ErrTeamExists   = errors.New("team already exists")
	ErrTeamNotFound = errors.New("team not found")

	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	// Pull Request errors
	ErrPRExists   = errors.New("pull request already exists")
	ErrPRNotFound = errors.New("pull request not found")
	ErrPRMerged   = errors.New("cannot modify merged pull request")

	// Reviewer errors
	ErrReviewerNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidates        = errors.New("no active candidates available for assignment")
)

// MapErrorToHTTPStatus maps domain errors to HTTP status codes
func MapErrorToHTTPStatus(err error) int {
	switch {
	case errors.Is(err, ErrTeamExists):
		return http.StatusBadRequest

	case errors.Is(err, ErrUserAlreadyExists),
		errors.Is(err, ErrPRExists),
		errors.Is(err, ErrPRMerged),
		errors.Is(err, ErrReviewerNotAssigned),
		errors.Is(err, ErrNoCandidates):
		return http.StatusConflict

	case errors.Is(err, ErrTeamNotFound),
		errors.Is(err, ErrUserNotFound),
		errors.Is(err, ErrPRNotFound):
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}

// MapErrorToErrorCode maps domain errors to API error codes
func MapErrorToErrorCode(err error) response.ErrorCode {
	switch {
	case errors.Is(err, ErrTeamExists):
		return response.ErrorCodeTeamExists
	case errors.Is(err, ErrPRExists):
		return response.ErrorCodePRExists
	case errors.Is(err, ErrPRMerged):
		return response.ErrorCodePRMerged
	case errors.Is(err, ErrReviewerNotAssigned):
		return response.ErrorCodeNotAssigned
	case errors.Is(err, ErrNoCandidates):
		return response.ErrorCodeNoCandidate
	case errors.Is(err, ErrTeamNotFound),
		errors.Is(err, ErrUserNotFound),
		errors.Is(err, ErrPRNotFound):
		return response.ErrorCodeNotFound
	default:
		return response.ErrorCodeNotFound
	}
}
