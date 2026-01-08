package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/myapp/internal/container"
	"example.com/myapp/internal/routes"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("Starting application")

	// Initialize container with all dependencies
	c := container.NewContainer()

	// Setup routes
	r := routes.SetupRoutes(c)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for errors from server
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", "addr", srv.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErrors:
		if err != http.ErrServerClosed {
			logger.Error("Server error", "error", err)
		}
	case sig := <-sigChan:
		logger.Info("Received signal, shutting down", "signal", sig)

		// Create a context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Shutdown the server
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Server shutdown error", "error", err)
			os.Exit(1)
		}

		logger.Info("Server shutdown complete")
	}
}