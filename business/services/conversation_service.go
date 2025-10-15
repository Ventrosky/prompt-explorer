package services

import (
	"fmt"

	"github.com/ventrosky/prompt-explorer/business/models"
)

// ConversationService handles conversation business logic
type ConversationService struct {
	// TODO: Add database repository when implemented
}

// NewConversationService creates a new ConversationService
func NewConversationService() *ConversationService {
	return &ConversationService{}
}

// CreateConversation creates a new conversation
func (s *ConversationService) CreateConversation(title string) (*models.Conversation, error) {
	// TODO: Implement database persistence
	conversation := models.NewConversation(title)
	return conversation, nil
}

// CreateConversationWithSystemPrompt creates a new conversation with an initial system message
func (s *ConversationService) CreateConversationWithSystemPrompt(title, systemPrompt string) (*models.Conversation, *models.Message, error) {
	// TODO: Implement database persistence
	conversation := models.NewConversation(title)
	systemMessage := models.NewMessage(conversation.ID, models.MessageTypeSystem, systemPrompt)
	return conversation, systemMessage, nil
}

// GetConversation retrieves a conversation by ID
func (s *ConversationService) GetConversation(id string) (*models.Conversation, error) {
	// TODO: Implement database retrieval
	return nil, fmt.Errorf("not implemented")
}

// ListConversations retrieves all conversations
func (s *ConversationService) ListConversations() ([]*models.Conversation, error) {
	// TODO: Implement database retrieval
	return []*models.Conversation{}, nil
}

// ToggleFavorite toggles the favorite status of a conversation
func (s *ConversationService) ToggleFavorite(id string) error {
	// TODO: Implement database update
	return fmt.Errorf("not implemented")
}

// GetSystemPrompt retrieves the system prompt (first system message) from a conversation
func (s *ConversationService) GetSystemPrompt(conversationID string) (string, error) {
	// TODO: Implement database retrieval
	// This would query for the first message with MessageTypeSystem in the conversation
	return "", fmt.Errorf("not implemented")
}
