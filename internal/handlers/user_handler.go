package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/myapp/internal/middleware"
	"example.com/myapp/internal/models"
	"example.com/myapp/internal/repositories"
	"example.com/myapp/internal/services"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service *services.UserService
	repo    repositories.UserRepository
}

func NewUserHandler(service *services.UserService, repo repositories.UserRepository) *UserHandler {
	return &UserHandler{
		service: service,
		repo:    repo,
	}
}

// RegisterRoutes registers all user-related routes with appropriate middleware
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		// Middleware for ALL user routes
		r.Use(middleware.LoggingMiddleware)
		r.Use(middleware.AuthMiddleware)

		r.Post("/", h.CreateUser)
		r.Get("/", h.GetAllUsers)

		// Nested route for ID-specific operations
		r.Route("/{id}", func(r chi.Router) {
			// Middleware ONLY for routes with {id}
			// Validates and extracts the ID to context
			// Loads the full user object into context
			r.Use(middleware.LoadUserMiddleware(h.repo))
			r.Use(middleware.ValidateIDMiddleware)

			r.Get("/", h.GetUser)
			r.Put("/", h.UpdateUser)
			r.Delete("/", h.DeleteUser)
		})
	})
}

// Create a new user - POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	
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
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
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
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context (loaded by LoadUserMiddleware)
	user, ok := r.Context().Value("user").(*models.User)
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Update a user - PUT /users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Extract validated ID from context (set by ValidateIDMiddleware)
	id, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	var req models.UpdateUserRequest
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
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
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
