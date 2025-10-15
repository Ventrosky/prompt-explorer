package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ventrosky/prompt-explorer/business/models"
)

// MessageRepository defines the interface for message data operations
type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error)
	GetByConversationID(ctx context.Context, conversationID uuid.UUID) ([]*models.Message, error)
	GetSystemPrompt(ctx context.Context, conversationID uuid.UUID) (string, error)
	Update(ctx context.Context, message *models.Message) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByConversationID(ctx context.Context, conversationID uuid.UUID) error
}

// messageRepository implements MessageRepository
type messageRepository struct {
	// TODO: Add database connection when implementing
}

// NewMessageRepository creates a new MessageRepository
func NewMessageRepository() MessageRepository {
	return &messageRepository{}
}

// Create creates a new message
func (r *messageRepository) Create(ctx context.Context, message *models.Message) error {
	// TODO: Implement database insert
	return fmt.Errorf("not implemented")
}

// GetByID retrieves a message by ID
func (r *messageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	// TODO: Implement database select
	return nil, fmt.Errorf("not implemented")
}

// GetByConversationID retrieves all messages for a conversation
func (r *messageRepository) GetByConversationID(ctx context.Context, conversationID uuid.UUID) ([]*models.Message, error) {
	// TODO: Implement database select
	return []*models.Message{}, nil
}

// GetSystemPrompt retrieves the system prompt (first system message) for a conversation
func (r *messageRepository) GetSystemPrompt(ctx context.Context, conversationID uuid.UUID) (string, error) {
	// TODO: Implement database select
	return "", fmt.Errorf("not implemented")
}

// Update updates a message
func (r *messageRepository) Update(ctx context.Context, message *models.Message) error {
	// TODO: Implement database update
	return fmt.Errorf("not implemented")
}

// Delete deletes a message
func (r *messageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// TODO: Implement database delete
	return fmt.Errorf("not implemented")
}

// DeleteByConversationID deletes all messages for a conversation
func (r *messageRepository) DeleteByConversationID(ctx context.Context, conversationID uuid.UUID) error {
	// TODO: Implement database delete
	return fmt.Errorf("not implemented")
}
