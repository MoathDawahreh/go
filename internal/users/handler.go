package users

import (
	"encoding/json"
	"net/http"

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
			// Validates and extracts the ID to context
			// Loads the full user object into context
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Get all users - GET /users
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// User is already fetched and in context by LoadUserMiddleware
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context (loaded by LoadUserMiddleware)
	user, ok := r.Context().Value("user").(*User)
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
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
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.UpdateUser(id, &req)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// We could also get the user from context if we wanted pre-fetched data

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Delete a user - DELETE /users/{id}
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract validated ID from context (set by ValidateIDMiddleware)
	id, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteUser(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
