package middleware

import (
	"context"
	"net/http"
)

// AuthMiddleware is a simple auth middleware that checks for authorization header
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Pass authenticated user info to context
		ctx := context.WithValue(r.Context(), "user", authHeader)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
