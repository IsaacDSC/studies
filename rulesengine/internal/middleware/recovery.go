package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"rulesengine/internal/shared"
)

// Recovery recovers from panics and returns 500 error
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v\n%s", err, debug.Stack())
				shared.JSON(w, http.StatusInternalServerError, map[string]string{
					"error":   "internal_error",
					"message": "an unexpected error occurred",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}


