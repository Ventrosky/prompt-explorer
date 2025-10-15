package api

import (
	"encoding/json"
	"net/http"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// Handlers contains all HTTP handlers
type Handlers struct{}

// NewHandlers creates a new Handlers instance
func NewHandlers() *Handlers {
	return &Handlers{}
}

// Health handles the health check endpoint
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{
		Status:  "ok",
		Service: "prompt-explorer",
	}

	json.NewEncoder(w).Encode(response)
}

// Home handles the root endpoint
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Prompt Explorer</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100 min-h-screen">
    <div class="container mx-auto px-4 py-8">
        <div class="max-w-4xl mx-auto">
            <header class="text-center mb-8">
                <h1 class="text-4xl font-bold text-gray-800 mb-4">Prompt Explorer</h1>
                <p class="text-lg text-gray-600">A tool for testing and managing system prompts with the Gemini API</p>
            </header>
            
            <div class="bg-white rounded-lg shadow-md p-6">
                <div class="text-center">
                    <h2 class="text-2xl font-semibold text-gray-700 mb-4">Welcome to Prompt Explorer</h2>
                    <p class="text-gray-600 mb-6">Create, test, and manage system prompts for your AI conversations.</p>
                    
                    <div class="space-y-4">
                        <a href="/api/health" class="inline-block bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded transition duration-200">
                            Health Check
                        </a>
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
