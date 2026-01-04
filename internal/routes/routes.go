package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"example.com/myapp/internal/handlers"
)

func SetupRoutes(userHandler *handlers.UserHandler, mediaHandler *handlers.MediaHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// User routes
	r.Post("/users", userHandler.CreateUser)
	r.Get("/users", userHandler.GetAllUsers)
	r.Get("/users/{id}", userHandler.GetUser)
	r.Put("/users/{id}", userHandler.UpdateUser)
	r.Delete("/users/{id}", userHandler.DeleteUser)

	// Media routes
	r.Post("/media/upload", mediaHandler.UploadMedia)
	r.Get("/media", mediaHandler.GetAllMedia)
	r.Get("/media/{id}", mediaHandler.GetMedia)
	r.Get("/media/{id}/download", mediaHandler.DownloadMedia)
	r.Delete("/media/{id}", mediaHandler.DeleteMedia)

	return r
}
