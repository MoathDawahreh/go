package users

import (
	"encoding/json"
	"log/slog"
	"net/http"

	appErr "example.com/myapp/internal/errors"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
	repo    Repository
}

func NewHandler(service *Service, repo Repository) *Handler {
	return &Handler{
		service: service,
		repo:    repo,
	}
}

// RegisterRoutes registers all user-related routes with appropriate middleware
// Middleware is passed as parameters to avoid circular imports
func (h *Handler) RegisterRoutes(r chi.Router, loggingMw, authMw, loadUserMw func(http.Handler) http.Handler, validateIDMw func(http.Handler) http.Handler) {
	r.Route("/users", func(r chi.Router) {
		// Middleware for ALL user routes
		r.Use(loggingMw)
		r.Use(authMw)

		r.Post("/", h.CreateUser)
		r.Get("/", h.GetAllUsers)

		// Nested route for ID-specific operations
		r.Route("/{id}", func(r chi.Router) {
			// Middleware ONLY for routes with {id}
			r.Use(loadUserMw)
			r.Use(validateIDMw)

			r.Get("/", h.GetUser)
			r.Put("/", h.UpdateUser)
			r.Delete("/", h.DeleteUser)
		})
	})
}

// Create a new user - POST /users
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request", "error", err)
		respondError(w, appErr.BadRequest("invalid request body"), http.StatusBadRequest)
		return
	}

	user, err := h.service.CreateUser(r.Context(), &req)
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		respondError(w, err, getStatusCode(err))
		return
	}

	slog.Info("User created", "user_id", user.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Get all users - GET /users
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		slog.Error("Failed to get all users", "error", err)
		respondError(w, err, getStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// Get a user - GET /users/{id}
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context (loaded by LoadUserMiddleware)
	user, ok := r.Context().Value("user").(*User)
	if !ok {
		slog.Error("User not found in context")
		respondError(w, appErr.NotFound("user not found"), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Update a user - PUT /users/{id}
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Extract validated ID from context (set by ValidateIDMiddleware)
	id, ok := r.Context().Value("userID").(int)
	if !ok {
		respondError(w, appErr.InvalidID("invalid user id"), http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request", "error", err)
		respondError(w, appErr.BadRequest("invalid request body"), http.StatusBadRequest)
		return
	}

	user, err := h.service.UpdateUser(r.Context(), id, &req)
	if err != nil {
		slog.Error("Failed to update user", "id", id, "error", err)
		respondError(w, err, getStatusCode(err))
		return
	}

	slog.Info("User updated", "user_id", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Delete a user - DELETE /users/{id}
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract validated ID from context (set by ValidateIDMiddleware)
	id, ok := r.Context().Value("userID").(int)
	if !ok {
		respondError(w, appErr.InvalidID("invalid user id"), http.StatusBadRequest)
		return
	}

	err := h.service.DeleteUser(r.Context(), id)
	if err != nil {
		slog.Error("Failed to delete user", "id", id, "error", err)
		respondError(w, err, getStatusCode(err))
		return
	}

	slog.Info("User deleted", "user_id", id)
	w.WriteHeader(http.StatusNoContent)
}

// Helper functions

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	ae := appErr.GetAppError(err)
	if ae == nil {
		return http.StatusInternalServerError
	}

	switch ae.Code {
	case appErr.ErrCodeNotFound:
		return http.StatusNotFound
	case appErr.ErrCodeBadRequest, appErr.ErrCodeInvalidID:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func respondError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"error": err.Error(),
	}
	json.NewEncoder(w).Encode(response)
}
