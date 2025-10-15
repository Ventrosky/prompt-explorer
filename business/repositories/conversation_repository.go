package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ventrosky/prompt-explorer/business/models"
)

// ConversationRepository defines the interface for conversation data operations
type ConversationRepository interface {
	Create(ctx context.Context, conversation *models.Conversation) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error)
	List(ctx context.Context) ([]*models.Conversation, error)
	Update(ctx context.Context, conversation *models.Conversation) error
	Delete(ctx context.Context, id uuid.UUID) error
	ToggleFavorite(ctx context.Context, id uuid.UUID) error
}

// conversationRepository implements ConversationRepository
type conversationRepository struct {
	// TODO: Add database connection when implementing
}

// NewConversationRepository creates a new ConversationRepository
func NewConversationRepository() ConversationRepository {
	return &conversationRepository{}
}

// Create creates a new conversation
func (r *conversationRepository) Create(ctx context.Context, conversation *models.Conversation) error {
	// TODO: Implement database insert
	return fmt.Errorf("not implemented")
}

// GetByID retrieves a conversation by ID
func (r *conversationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error) {
	// TODO: Implement database select
	return nil, fmt.Errorf("not implemented")
}

// List retrieves all conversations
func (r *conversationRepository) List(ctx context.Context) ([]*models.Conversation, error) {
	// TODO: Implement database select
	return []*models.Conversation{}, nil
}

// Update updates a conversation
func (r *conversationRepository) Update(ctx context.Context, conversation *models.Conversation) error {
	// TODO: Implement database update
	return fmt.Errorf("not implemented")
}

// Delete deletes a conversation
func (r *conversationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// TODO: Implement database delete
	return fmt.Errorf("not implemented")
}

// ToggleFavorite toggles the favorite status of a conversation
func (r *conversationRepository) ToggleFavorite(ctx context.Context, id uuid.UUID) error {
	// TODO: Implement database update
	return fmt.Errorf("not implemented")
}
