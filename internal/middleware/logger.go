package middleware

import (
	"net/http"
	"time"

	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

// responseWriter is a wrapper around http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// newResponseWriter creates a new responseWriter
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // default status code
		written:        false,
	}
}

// WriteHeader captures the status code and calls the underlying WriteHeader
func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.written {
		rw.statusCode = statusCode
		rw.written = true
		rw.ResponseWriter.WriteHeader(statusCode)
	}
}

// Write ensures WriteHeader is called and delegates to the underlying Write
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// Logger is a middleware that logs HTTP requests
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status code
		wrappedWriter := newResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(wrappedWriter, r)

		// Log the request details
		duration := time.Since(start)
		logger.Info(
			"HTTP %s %s - Status: %d - Duration: %v",
			r.Method,
			r.URL.Path,
			wrappedWriter.statusCode,
			duration,
		)
	})
}
