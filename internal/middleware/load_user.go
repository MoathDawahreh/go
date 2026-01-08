package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"example.com/myapp/internal/users"
)

// LoadUserMiddleware fetches the user from the database and puts it into the request context.
// This middleware assumes the {id} parameter has been validated by ValidateIDMiddleware
// Usage: Apply after ValidateIDMiddleware
func LoadUserMiddleware(repo users.Repository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get validated ID from context (set by ValidateIDMiddleware)
			id, ok := r.Context().Value("userID").(int)
			if !ok {
				slog.Error("User ID not found in context")
				http.Error(w, "User ID not found in context", http.StatusInternalServerError)
				return
			}

			// Fetch user from repository with request context
			user, err := repo.GetByID(r.Context(), id)
			if err != nil {
				slog.Error("Failed to load user", "id", id, "error", err)
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			// Put user in context for handler to use
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
