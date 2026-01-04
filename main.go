package main

import (
	"fmt"
	"net/http"

	"example.com/myapp/internal/container"
	"example.com/myapp/internal/routes"
)

func main() {
	// Initialize container with all dependencies
	c := container.NewContainer()
	
	// Setup routes
	r := routes.SetupRoutes(c.UserHandler, c.MediaHandler)
	
	// Start server
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}