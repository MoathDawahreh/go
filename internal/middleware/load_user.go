package middleware

import (
	"context"
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
				http.Error(w, "User ID not found in context", http.StatusInternalServerError)
				return
			}

			// Fetch user from repository
			user, err := repo.GetByID(id)
			if err != nil {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			// Put user in context for handler to use
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
