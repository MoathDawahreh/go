package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"example.com/myapp/internal/container"
	mw "example.com/myapp/internal/middleware"
)

func SetupRoutes(c *container.Container) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Register handler routes with middleware
	c.UserHandler.RegisterRoutes(r, mw.LoggingMiddleware, mw.AuthMiddleware, mw.LoadUserMiddleware(c.UserRepository), mw.ValidateIDMiddleware)
	c.MediaHandler.RegisterRoutes(r, mw.LoggingMiddleware, mw.AuthMiddleware, mw.ValidateIDMiddleware)

	return r
}
