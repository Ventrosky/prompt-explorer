package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ventrosky/prompt-explorer/business/services"
)

// Routes sets up all the HTTP routes
func Routes(conversationService *services.ConversationService) http.Handler {
	r := mux.NewRouter()

	handlers := NewHandlers(conversationService)

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", handlers.Health).Methods("GET")

	// Conversation routes
	api.HandleFunc("/conversations", handlers.CreateConversation).Methods("POST")
	api.HandleFunc("/conversations/{id}", handlers.GetConversation).Methods("GET")
	api.HandleFunc("/conversations/{id}/messages", handlers.SendMessage).Methods("POST")
	api.HandleFunc("/conversations/{id}/favorite", handlers.ToggleFavorite).Methods("POST")

	// Simple health check interface
	r.HandleFunc("/", handlers.Home).Methods("GET")

	return r
}
