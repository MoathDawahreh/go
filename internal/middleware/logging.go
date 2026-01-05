package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// LoggingMiddleware logs HTTP request details
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request
		fmt.Printf("[%s] %s %s\n", r.Method, r.RequestURI, time.Now().Format(time.RFC3339))

		next.ServeHTTP(w, r)

		// Log response time
		duration := time.Since(start)
		fmt.Printf("  └─ completed in %v\n", duration)
	})
}
