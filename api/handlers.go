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
</head>
<body>
    <h1>Prompt Explorer</h1>
    <p>Welcome to Prompt Explorer - A tool for testing and managing system prompts with the Gemini API.</p>
    <p>Status: <a href="/health">Health Check</a></p>
</body>
</html>
	`

	w.Write([]byte(html))
}
