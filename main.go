package main

import (
	"fmt"
	"net/http"

	"example.com/myapp/internal/handlers"
	"example.com/myapp/internal/routes"
	"example.com/myapp/internal/services"
)

func main() {
	// Initialize services
	userService := services.NewUserService()
	mediaService := services.NewMediaService()
	
	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	mediaHandler := handlers.NewMediaHandler(mediaService)
	
	// Setup routes
	r := routes.SetupRoutes(userHandler, mediaHandler)
	
	// Start server
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}