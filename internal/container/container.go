package container

import (
	"example.com/myapp/internal/handlers"
	"example.com/myapp/internal/repositories"
	"example.com/myapp/internal/services"
)

type Container struct {
	// Repositories
	UserRepository  repositories.UserRepository
	MediaRepository repositories.MediaRepository

	// Services
	UserService   *services.UserService
	MediaService  *services.MediaService

	// Handlers
	UserHandler   *handlers.UserHandler
	MediaHandler  *handlers.MediaHandler
}

func NewContainer() *Container {
	// Initialize repositories
	userRepo := repositories.NewInMemoryUserRepository()
	mediaRepo := repositories.NewInMemoryMediaRepository()

	// Initialize services with repositories
	userService := services.NewUserService(userRepo)
	mediaService := services.NewMediaService(mediaRepo)

	// Initialize handlers with services and repositories
	userHandler := handlers.NewUserHandler(userService, userRepo)
	mediaHandler := handlers.NewMediaHandler(mediaService)

	return &Container{
		UserRepository:  userRepo,
		MediaRepository: mediaRepo,
		UserService:     userService,
		MediaService:    mediaService,
		UserHandler:     userHandler,
		MediaHandler:    mediaHandler,
	}
}