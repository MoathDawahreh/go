package container

import (
	"example.com/myapp/internal/media"
	"example.com/myapp/internal/users"
)

type Container struct {
	// Repositories
	UserRepository  users.Repository
	MediaRepository media.Repository

	// Services
	UserService   *users.Service
	MediaService  *media.Service

	// Handlers
	UserHandler   *users.Handler
	MediaHandler  *media.Handler
}

func NewContainer() *Container {
	// Initialize repositories
	userRepo := users.NewInMemoryRepository()
	mediaRepo := media.NewInMemoryRepository()

	// Initialize services with repositories
	userService := users.NewService(userRepo)
	mediaService := media.NewService(mediaRepo)

	// Initialize handlers with services and repositories
	userHandler := users.NewHandler(userService, userRepo)
	mediaHandler := media.NewHandler(mediaService)

	return &Container{
		UserRepository:  userRepo,
		MediaRepository: mediaRepo,
		UserService:     userService,
		MediaService:    mediaService,
		UserHandler:     userHandler,
		MediaHandler:    mediaHandler,
	}
}