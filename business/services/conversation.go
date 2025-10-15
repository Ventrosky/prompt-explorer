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
	conversationID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation ID: %w", err)
	}

	return s.conversationRepo.GetByID(ctx, conversationID)
}

// ListConversations retrieves all conversations
func (s *ConversationService) ListConversations(ctx context.Context) ([]*models.Conversation, error) {
	return s.conversationRepo.List(ctx)
}

// ToggleFavorite toggles the favorite status of a conversation
func (s *ConversationService) ToggleFavorite(ctx context.Context, id string) error {
	conversationID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid conversation ID: %w", err)
	}

	return s.conversationRepo.ToggleFavorite(ctx, conversationID)
}

// GetSystemPrompt retrieves the system prompt (first system message) from a conversation
func (s *ConversationService) GetSystemPrompt(ctx context.Context, conversationID string) (string, error) {
	convID, err := uuid.Parse(conversationID)
	if err != nil {
		return "", fmt.Errorf("invalid conversation ID: %w", err)
	}

	// Get the first system message for this conversation
	messages, err := s.messageRepo.GetByConversationID(ctx, convID)
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

	return s.conversationRepo.GetNextConversation(ctx, currentUUID)
}

// GetPreviousConversation gets the previous conversation (older) by timestamp
func (s *ConversationService) GetPreviousConversation(ctx context.Context, currentID string) (*models.Conversation, error) {
	currentUUID, err := uuid.Parse(currentID)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation ID: %w", err)
	}

	return s.conversationRepo.GetPreviousConversation(ctx, currentUUID)
}

// GetConversationHistory gets a conversation with all its messages
func (s *ConversationService) GetConversationHistory(ctx context.Context, conversationID string) (*models.Conversation, []*models.Message, error) {
	convID, err := uuid.Parse(conversationID)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid conversation ID: %w", err)
	}

	// Get conversation
	conversation, err := s.conversationRepo.GetByID(ctx, convID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Get messages
	messages, err := s.messageRepo.GetByConversationID(ctx, convID)
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
	_, err = s.conversationRepo.GetByID(ctx, convID)
	if err != nil {
		return nil, fmt.Errorf("conversation not found: %w", err)
	}

	// Create message
	message := models.NewMessage(convID, sender, content)

	// Save message
	if err := s.messageRepo.Create(ctx, message); err != nil {
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
	message, err := s.messageRepo.GetByID(ctx, msgID)
	if err != nil {
		return fmt.Errorf("message not found: %w", err)
	}

	// Update content
	message.Content = content

	// Save changes
	if err := s.messageRepo.Update(ctx, message); err != nil {
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
	if err := s.messageRepo.Delete(ctx, msgID); err != nil {
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

	message, err := s.messageRepo.GetByID(ctx, msgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return message, nil
}
