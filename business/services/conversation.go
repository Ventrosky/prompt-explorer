package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ventrosky/prompt-explorer/business/models"
	"github.com/ventrosky/prompt-explorer/business/repositories"
)

// ConversationService handles conversation business logic
type ConversationService struct {
	repo          repositories.Repository
	geminiService *GeminiService
}

// NewConversationService creates a new ConversationService
func NewConversationService(repo repositories.Repository, geminiService *GeminiService) *ConversationService {
	return &ConversationService{
		repo:          repo,
		geminiService: geminiService,
	}
}

// CreateConversation creates a new conversation
func (s *ConversationService) CreateConversation(ctx context.Context, title string) (*models.Conversation, error) {
	conversation := models.NewConversation(title)

	if err := s.repo.CreateConversation(ctx, conversation); err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return conversation, nil
}

// CreateConversationWithSystemPrompt creates a new conversation with an initial system message
func (s *ConversationService) CreateConversationWithSystemPrompt(ctx context.Context, title, systemPrompt string) (*models.Conversation, *models.Message, error) {
	conversation := models.NewConversation(title)

	// Create conversation first
	if err := s.repo.CreateConversation(ctx, conversation); err != nil {
		return nil, nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Create system message
	systemMessage := models.NewMessage(conversation.ID, models.MessageTypeSystem, systemPrompt)
	if err := s.repo.CreateMessage(ctx, systemMessage); err != nil {
		return nil, nil, fmt.Errorf("failed to create system message: %w", err)
	}

	return conversation, systemMessage, nil
}

// GetConversation retrieves a conversation by ID
func (s *ConversationService) GetConversation(ctx context.Context, id string) (*models.Conversation, error) {
	conversationID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation ID: %w", err)
	}

	return s.repo.GetConversationByID(ctx, conversationID)
}

// ListConversations retrieves all conversations
func (s *ConversationService) ListConversations(ctx context.Context) ([]*models.Conversation, error) {
	return s.repo.ListConversations(ctx)
}

// ToggleFavorite toggles the favorite status of a conversation
func (s *ConversationService) ToggleFavorite(ctx context.Context, id string) error {
	conversationID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	return s.repo.ToggleConversationFavorite(ctx, conversationID)
}

// GetSystemPrompt retrieves the system prompt (first system message) from a conversation
func (s *ConversationService) GetSystemPrompt(ctx context.Context, conversationID string) (string, error) {
	convID, err := uuid.Parse(conversationID)
	if err != nil {
		return "", fmt.Errorf("invalid conversation ID: %w", err)
	}

	// Get the first system message for this conversation
	messages, err := s.repo.GetMessagesByConversationID(ctx, convID)
	if err != nil {
		return "", fmt.Errorf("failed to get messages: %w", err)
	}

	// Find the first system message
	for _, msg := range messages {
		if msg.Sender == models.MessageTypeSystem {
			return msg.Content, nil
		}
	}

	return "", fmt.Errorf("no system prompt found for conversation")
}

// GetNextConversation gets the next conversation (newer) by timestamp
func (s *ConversationService) GetNextConversation(ctx context.Context, currentID string) (*models.Conversation, error) {
	currentUUID, err := uuid.Parse(currentID)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation ID: %w", err)
	}

	return s.repo.GetNextConversation(ctx, currentUUID)
}

// GetPreviousConversation gets the previous conversation (older) by timestamp
func (s *ConversationService) GetPreviousConversation(ctx context.Context, currentID string) (*models.Conversation, error) {
	currentUUID, err := uuid.Parse(currentID)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation ID: %w", err)
	}

	return s.repo.GetPreviousConversation(ctx, currentUUID)
}

// GetConversationHistory gets a conversation with all its messages
func (s *ConversationService) GetConversationHistory(ctx context.Context, conversationID string) (*models.Conversation, []*models.Message, error) {
	convID, err := uuid.Parse(conversationID)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid conversation ID: %w", err)
	}

	// Get conversation
	conversation, err := s.repo.GetConversationByID(ctx, convID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Get messages
	messages, err := s.repo.GetMessagesByConversationID(ctx, convID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return conversation, messages, nil
}

// AddMessage adds a new message to a conversation
func (s *ConversationService) AddMessage(ctx context.Context, conversationID string, sender models.MessageType, content string) (*models.Message, error) {
	convID, err := uuid.Parse(conversationID)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation ID: %w", err)
	}

	// Verify conversation exists
	_, err = s.repo.GetConversationByID(ctx, convID)
	if err != nil {
		return nil, fmt.Errorf("conversation not found: %w", err)
	}

	// Create message
	message := models.NewMessage(convID, sender, content)

	// Save message
	if err := s.repo.CreateMessage(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return message, nil
}

// UpdateMessage updates an existing message
func (s *ConversationService) UpdateMessage(ctx context.Context, messageID string, content string) error {
	msgID, err := uuid.Parse(messageID)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	// Get existing message
	message, err := s.repo.GetMessageByID(ctx, msgID)
	if err != nil {
		return fmt.Errorf("message not found: %w", err)
	}

	// Update content
	message.Content = content

	// Save changes
	if err := s.repo.UpdateMessage(ctx, message); err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	return nil
}

// DeleteMessage deletes a message
func (s *ConversationService) DeleteMessage(ctx context.Context, messageID string) error {
	msgID, err := uuid.Parse(messageID)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	// Delete message
	if err := s.repo.DeleteMessage(ctx, msgID); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

// GetMessage gets a specific message by ID
func (s *ConversationService) GetMessage(ctx context.Context, messageID string) (*models.Message, error) {
	msgID, err := uuid.Parse(messageID)
	if err != nil {
		return nil, fmt.Errorf("invalid message ID: %w", err)
	}

	message, err := s.repo.GetMessageByID(ctx, msgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return message, nil
}

// SendMessageToGemini sends a user message to Gemini and returns the assistant response
func (s *ConversationService) SendMessageToGemini(ctx context.Context, conversationID string, userMessage string) (*models.Message, error) {
	// Check if Gemini service is available
	if s.geminiService == nil {
		return nil, fmt.Errorf("Gemini service not available")
	}

	convID, err := uuid.Parse(conversationID)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation ID: %w", err)
	}

	// Get conversation to verify it exists
	_, err = s.repo.GetConversationByID(ctx, convID)
	if err != nil {
		return nil, fmt.Errorf("conversation not found: %w", err)
	}

	// Add user message to conversation
	_, err = s.AddMessage(ctx, conversationID, models.MessageTypeUser, userMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to add user message: %w", err)
	}

	// Get conversation history for context
	_, messages, err := s.GetConversationHistory(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation history: %w", err)
	}

	// Get system prompt (first system message)
	systemPrompt, err := s.repo.GetSystemPrompt(ctx, convID)
	if err != nil {
		// If no system prompt, use a default one
		systemPrompt = "You are a helpful AI assistant. Please respond to the user's messages in a helpful and informative way."
	}

	// Ensure system prompt is the first message
	var allMessages []*models.Message
	if len(messages) == 0 || messages[0].Sender != models.MessageTypeSystem {
		// Add system prompt as first message if not already present
		systemMsg := models.NewMessage(convID, models.MessageTypeSystem, systemPrompt)
		allMessages = append([]*models.Message{systemMsg}, messages...)
	} else {
		allMessages = messages
	}

	// Send to Gemini
	response, err := s.geminiService.SendMessage(ctx, conversationID, allMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to send message to Gemini: %w", err)
	}

	// Add assistant response to conversation
	assistantMsg, err := s.AddMessage(ctx, conversationID, models.MessageTypeAssistant, response)
	if err != nil {
		return nil, fmt.Errorf("failed to add assistant message: %w", err)
	}

	return assistantMsg, nil
}

// HealthCheck checks if Gemini API is accessible
func (s *ConversationService) HealthCheck(ctx context.Context) error {
	if s.geminiService == nil {
		return fmt.Errorf("Gemini service not available")
	}

	return s.geminiService.HealthCheck(ctx)
}
