package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware logs HTTP request details with structured logging
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request
		slog.Info("Request received",
			"method", r.Method,
			"path", r.RequestURI,
			"remote_addr", r.RemoteAddr,
		)

		next.ServeHTTP(w, r)

		// Log response time
		duration := time.Since(start)
		slog.Info("Request completed",
			"method", r.Method,
			"path", r.RequestURI,
			"duration_ms", duration.Milliseconds(),
		)
	})
}

