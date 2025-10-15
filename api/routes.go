package api

import (
	"net/http"
)

// Routes sets up all the HTTP routes
func Routes() *http.ServeMux {
	mux := http.NewServeMux()
	handlers := NewHandlers()

	// Health check endpoint
	mux.HandleFunc("/health", handlers.Health)

	// Root endpoint
	mux.HandleFunc("/", handlers.Home)

	return mux
}
