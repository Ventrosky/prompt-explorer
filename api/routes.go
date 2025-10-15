package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Routes sets up all the HTTP routes
func Routes() http.Handler {
	r := mux.NewRouter()

	handlers := NewHandlers()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", handlers.Health).Methods("GET")

	// Root endpoint
	r.HandleFunc("/", handlers.Home).Methods("GET")

	return r
}
