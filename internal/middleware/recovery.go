package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

// Recovery is a middleware that recovers from panics and returns 500 Internal Server Error
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				stackTrace := debug.Stack()
				logger.Error(
					"PANIC recovered: %v\nStack trace:\n%s",
					err,
					string(stackTrace),
				)

				// Return 500 Internal Server Error
				http.Error(
					w,
					fmt.Sprintf("Internal Server Error: %v", err),
					http.StatusInternalServerError,
				)
			}
		}()

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
