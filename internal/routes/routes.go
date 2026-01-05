package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"example.com/myapp/internal/container"
)

func SetupRoutes(c *container.Container) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Register handler routes
	c.UserHandler.RegisterRoutes(r)
	c.MediaHandler.RegisterRoutes(r)

	return r
}
