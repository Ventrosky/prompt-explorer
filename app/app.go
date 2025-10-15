package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ventrosky/prompt-explorer/api"
	"github.com/ventrosky/prompt-explorer/foundation/config"
	"github.com/ventrosky/prompt-explorer/foundation/database"
)

// App represents the application
type App struct {
	server   *http.Server
	config   *config.Config
	database *database.Database
}

// New creates a new App instance
func New() *App {
	cfg := config.Load()
	return &App{
		config:   cfg,
		database: database.New(cfg),
	}
}

// Start starts the application server
func (a *App) Start() error {
	// Connect to database
	if err := a.database.Connect(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create HTTP server
	a.server = &http.Server{
		Addr:    ":" + a.config.ServerPort,
		Handler: api.Routes(),
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Prompt Explorer server starting on port %s", a.config.ServerPort)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited")
	return nil
}

// Stop gracefully stops the application
func (a *App) Stop() error {
	if a.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown database connection
	if a.database != nil {
		a.database.Close()
	}

	return a.server.Shutdown(ctx)
}
