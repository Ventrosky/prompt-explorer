package services

import (
	"context"
	"fmt"
	"time"

	"github.com/ventrosky/prompt-explorer/business/models"
	"google.golang.org/genai"
)

// GeminiService handles interactions with the Gemini API
type GeminiService struct {
	client *genai.Client
	model  string
}

// NewGeminiService creates a new GeminiService
func NewGeminiService(apiKey string) (*GeminiService, error) {
	ctx := context.Background()

	// Create client with API key
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Use Gemini 2.5 Pro model (most capable model)
	modelName := "gemini-2.5-pro"

	return &GeminiService{
		client: client,
		model:  modelName,
	}, nil
}

// SendMessage sends a message to Gemini and returns the response
func (s *GeminiService) SendMessage(ctx context.Context, messages []*models.Message) (string, error) {
	// Prepare the conversation history
	var history []*genai.Content

	// Add conversation messages (first message should be system prompt)
	for _, msg := range messages {
		var role genai.Role
		switch msg.Sender {
		case models.MessageTypeSystem:
			role = genai.RoleUser // System messages are treated as user messages
		case models.MessageTypeUser:
			role = genai.RoleUser
		case models.MessageTypeAssistant:
			role = genai.RoleModel
		default:
			role = genai.RoleUser
		}

		history = append(history, genai.NewContentFromText(msg.Content, role))
	}

	// Generate content with timeout
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Create chat session with history
	chat, err := s.client.Chats.Create(ctx, s.model, nil, history)
	if err != nil {
		return "", fmt.Errorf("failed to create chat: %w", err)
	}

	// Send the last message (which should be the user's current message)
	if len(messages) == 0 {
		return "", fmt.Errorf("no messages to send")
	}

	lastMessage := messages[len(messages)-1]
	response, err := chat.SendMessage(ctx, genai.Part{Text: lastMessage.Content})
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	// Extract the response text
	if len(response.Candidates) == 0 {
		return "", fmt.Errorf("no response candidates received")
	}

	if len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response parts received")
	}

	// Get text from the first part
	part := response.Candidates[0].Content.Parts[0]
	if part.Text == "" {
		return "", fmt.Errorf("empty response received")
	}

	return part.Text, nil
}

// HealthCheck checks if the Gemini API is accessible
func (s *GeminiService) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Send a simple test message
	chat, err := s.client.Chats.Create(ctx, s.model, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create chat: %w", err)
	}

	response, err := chat.SendMessage(ctx, genai.Part{Text: "Hello, are you working?"})
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	if len(response.Candidates) == 0 {
		return fmt.Errorf("health check failed: no response candidates")
	}

	return nil
}
