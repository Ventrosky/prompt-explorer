package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ventrosky/prompt-explorer/business/models"
	"github.com/ventrosky/prompt-explorer/business/services"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Database  string `json:"database"`
	GeminiAPI string `json:"gemini_api"`
}

// Handlers contains all HTTP handlers
type Handlers struct {
	conversationService *services.ConversationService
}

// NewHandlers creates a new Handlers instance
func NewHandlers(conversationService *services.ConversationService) *Handlers {
	return &Handlers{
		conversationService: conversationService,
	}
}

// Health handles the health check endpoint
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	response := HealthResponse{
		Status:    "ok",
		Service:   "prompt-explorer",
		Database:  "ok",
		GeminiAPI: "unknown",
	}

	// Check Gemini API health if service is available
	if h.conversationService != nil {
		if err := h.conversationService.HealthCheck(ctx); err != nil {
			response.Status = "degraded"
			response.GeminiAPI = "error: " + err.Error()
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			response.GeminiAPI = "ok"
			w.WriteHeader(http.StatusOK)
		}
	} else {
		response.GeminiAPI = "not configured"
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}

// Home handles the root endpoint - simple health check interface
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Prompt Explorer - Backend</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 min-h-screen">
    <div class="container mx-auto px-4 py-8">
        <div class="max-w-4xl mx-auto">
            <header class="text-center mb-8">
                <h1 class="text-4xl font-bold text-gray-800 mb-4">Prompt Explorer</h1>
                <p class="text-lg text-gray-600">Backend API Server</p>
            </header>
            
            <div class="bg-white rounded-lg shadow-md p-6">
                <div class="text-center">
                    <h2 class="text-2xl font-semibold text-gray-700 mb-4">Backend is Running</h2>
                    <p class="text-gray-600 mb-6">The Go backend API is ready to serve the React frontend.</p>
                    
                    <div class="space-y-4">
                        <a href="/api/health" class="inline-block bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded transition duration-200">
                            Health Check
                        </a>
                        <div class="text-sm text-gray-500">
                            <p>Frontend: <a href="http://localhost:3000" class="text-blue-500 hover:underline">http://localhost:3000</a></p>
                            <p>API: <a href="/api/health" class="text-blue-500 hover:underline">/api/health</a></p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</body>
</html>
	`

	w.Write([]byte(html))
}

// CreateConversation creates a new conversation
func (h *Handlers) CreateConversation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request body
	var req struct {
		Title        string `json:"title"`
		SystemPrompt string `json:"system_prompt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create conversation with system prompt
	conversation, systemMessage, err := h.conversationService.CreateConversationWithSystemPrompt(ctx, req.Title, req.SystemPrompt)
	if err != nil {
		http.Error(w, "Failed to create conversation", http.StatusInternalServerError)
		return
	}

	// Return conversation data
	response := struct {
		Conversation  *models.Conversation `json:"conversation"`
		SystemMessage *models.Message      `json:"system_message"`
	}{
		Conversation:  conversation,
		SystemMessage: systemMessage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetConversation retrieves a conversation by ID
func (h *Handlers) GetConversation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get conversation ID from URL path
	conversationID := r.URL.Path[len("/api/conversations/"):]
	if conversationID == "" {
		http.Error(w, "Conversation ID required", http.StatusBadRequest)
		return
	}

	// Get conversation
	conversation, err := h.conversationService.GetConversation(ctx, conversationID)
	if err != nil {
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	// Get conversation history
	_, messages, err := h.conversationService.GetConversationHistory(ctx, conversationID)
	if err != nil {
		http.Error(w, "Failed to load conversation history", http.StatusInternalServerError)
		return
	}

	// Get system prompt
	systemPrompt, err := h.conversationService.GetSystemPrompt(ctx, conversationID)
	if err != nil {
		systemPrompt = ""
	}

	// Prepare response
	response := struct {
		Conversation *models.Conversation `json:"conversation"`
		Messages     []*models.Message    `json:"messages"`
		SystemPrompt string               `json:"system_prompt"`
	}{
		Conversation: conversation,
		Messages:     messages,
		SystemPrompt: systemPrompt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SendMessage sends a message to a conversation
func (h *Handlers) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get conversation ID from URL path
	conversationID := r.URL.Path[len("/api/conversations/"):]
	conversationID = conversationID[:len(conversationID)-len("/messages")]
	if conversationID == "" {
		http.Error(w, "Conversation ID required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Send message to Gemini
	assistantMessage, err := h.conversationService.SendMessageToGemini(ctx, conversationID, req.Content)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	// Return assistant message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assistantMessage)
}

// ToggleFavorite toggles the favorite status of a conversation
func (h *Handlers) ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get conversation ID from URL path
	conversationID := r.URL.Path[len("/api/conversations/"):]
	conversationID = conversationID[:len(conversationID)-len("/favorite")]
	if conversationID == "" {
		http.Error(w, "Conversation ID required", http.StatusBadRequest)
		return
	}

	// Toggle favorite
	if err := h.conversationService.ToggleFavorite(ctx, conversationID); err != nil {
		http.Error(w, "Failed to toggle favorite", http.StatusInternalServerError)
		return
	}

	// Get updated conversation
	conversation, err := h.conversationService.GetConversation(ctx, conversationID)
	if err != nil {
		http.Error(w, "Failed to get updated conversation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}
