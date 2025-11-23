package handler

import (
	"encoding/json"
	"net/http"

	"avito-backend-trainee-assignment-autumn-2025/internal/dto/response"
	pkgerrors "avito-backend-trainee-assignment-autumn-2025/pkg/errors"
	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

// respondWithJSON sends a JSON response with the given status code
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			logger.Error("Failed to encode JSON response: %v", err)
		}
	}
}

// respondWithError sends an error response in the standard API format
func respondWithError(w http.ResponseWriter, err error) {
	statusCode := pkgerrors.MapErrorToHTTPStatus(err)
	errorCode := pkgerrors.MapErrorToErrorCode(err)

	errorResponse := response.ErrorResponse{
		Error: response.ErrorDetail{
			Code:    errorCode,
			Message: err.Error(),
		},
	}

	logger.Debug("Responding with error: status=%d, code=%s, message=%s", statusCode, errorCode, err.Error())
	respondWithJSON(w, statusCode, errorResponse)
}

// decodeJSONBody decodes JSON request body into the given struct
func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		logger.Error("Failed to decode JSON body: %v", err)
		return err
	}

	return nil
}
