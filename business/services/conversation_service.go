package services

import (
	"context"
	"fmt"

	"github.com/ventrosky/prompt-explorer/business/models"
	"github.com/ventrosky/prompt-explorer/business/repositories"
)

// ConversationService handles conversation business logic
type ConversationService struct {
	conversationRepo repositories.ConversationRepository
	messageRepo      repositories.MessageRepository
}

// NewConversationService creates a new ConversationService
func NewConversationService(conversationRepo repositories.ConversationRepository, messageRepo repositories.MessageRepository) *ConversationService {
	return &ConversationService{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
	}
}

// CreateConversation creates a new conversation
func (s *ConversationService) CreateConversation(ctx context.Context, title string) (*models.Conversation, error) {
	conversation := models.NewConversation(title)

	if err := s.conversationRepo.Create(ctx, conversation); err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return conversation, nil
}

// CreateConversationWithSystemPrompt creates a new conversation with an initial system message
func (s *ConversationService) CreateConversationWithSystemPrompt(ctx context.Context, title, systemPrompt string) (*models.Conversation, *models.Message, error) {
	conversation := models.NewConversation(title)

	// Create conversation first
	if err := s.conversationRepo.Create(ctx, conversation); err != nil {
		return nil, nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	// Create system message
	systemMessage := models.NewMessage(conversation.ID, models.MessageTypeSystem, systemPrompt)
	if err := s.messageRepo.Create(ctx, systemMessage); err != nil {
		return nil, nil, fmt.Errorf("failed to create system message: %w", err)
	}

	return conversation, systemMessage, nil
}

// GetConversation retrieves a conversation by ID
func (s *ConversationService) GetConversation(ctx context.Context, id string) (*models.Conversation, error) {
	// TODO: Parse UUID from string
	return nil, fmt.Errorf("not implemented")
}

// ListConversations retrieves all conversations
func (s *ConversationService) ListConversations(ctx context.Context) ([]*models.Conversation, error) {
	return s.conversationRepo.List(ctx)
}

// ToggleFavorite toggles the favorite status of a conversation
func (s *ConversationService) ToggleFavorite(ctx context.Context, id string) error {
	// TODO: Parse UUID from string
	return fmt.Errorf("not implemented")
}

// GetSystemPrompt retrieves the system prompt (first system message) from a conversation
func (s *ConversationService) GetSystemPrompt(ctx context.Context, conversationID string) (string, error) {
	// TODO: Parse UUID from string
	return "", fmt.Errorf("not implemented")
}
