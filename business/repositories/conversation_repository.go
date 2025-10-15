package repositories

import (
	"context"
	"database/sql"
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
	GetNextConversation(ctx context.Context, currentID uuid.UUID) (*models.Conversation, error)
	GetPreviousConversation(ctx context.Context, currentID uuid.UUID) (*models.Conversation, error)
}

// conversationRepository implements ConversationRepository
type conversationRepository struct {
	db *sql.DB
}

// NewConversationRepository creates a new ConversationRepository
func NewConversationRepository(db *sql.DB) ConversationRepository {
	return &conversationRepository{
		db: db,
	}
}

// Create creates a new conversation
func (r *conversationRepository) Create(ctx context.Context, conversation *models.Conversation) error {
	query := `
		INSERT INTO conversations (id, title, is_favorite, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query,
		conversation.ID,
		conversation.Title,
		conversation.IsFavorite,
		conversation.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	return nil
}

// GetByID retrieves a conversation by ID
func (r *conversationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error) {
	query := `
		SELECT id, title, is_favorite, created_at
		FROM conversations
		WHERE id = $1
	`

	var conversation models.Conversation
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&conversation.ID,
		&conversation.Title,
		&conversation.IsFavorite,
		&conversation.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return &conversation, nil
}

// List retrieves all conversations
func (r *conversationRepository) List(ctx context.Context) ([]*models.Conversation, error) {
	query := `
		SELECT id, title, is_favorite, created_at
		FROM conversations
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer rows.Close()

	var conversations []*models.Conversation
	for rows.Next() {
		var conversation models.Conversation
		err := rows.Scan(
			&conversation.ID,
			&conversation.Title,
			&conversation.IsFavorite,
			&conversation.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, &conversation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate conversations: %w", err)
	}

	return conversations, nil
}

// Update updates a conversation
func (r *conversationRepository) Update(ctx context.Context, conversation *models.Conversation) error {
	query := `
		UPDATE conversations
		SET title = $2, is_favorite = $3
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		conversation.ID,
		conversation.Title,
		conversation.IsFavorite,
	)

	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

// Delete deletes a conversation
func (r *conversationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM conversations WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

// ToggleFavorite toggles the favorite status of a conversation
func (r *conversationRepository) ToggleFavorite(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE conversations
		SET is_favorite = NOT is_favorite
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to toggle favorite: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

// GetNextConversation gets the next conversation (newer) by timestamp
func (r *conversationRepository) GetNextConversation(ctx context.Context, currentID uuid.UUID) (*models.Conversation, error) {
	query := `
		SELECT id, title, is_favorite, created_at
		FROM conversations
		WHERE created_at > (SELECT created_at FROM conversations WHERE id = $1)
		ORDER BY created_at ASC
		LIMIT 1
	`

	var conversation models.Conversation
	err := r.db.QueryRowContext(ctx, query, currentID).Scan(
		&conversation.ID,
		&conversation.Title,
		&conversation.IsFavorite,
		&conversation.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no next conversation")
		}
		return nil, fmt.Errorf("failed to get next conversation: %w", err)
	}

	return &conversation, nil
}

// GetPreviousConversation gets the previous conversation (older) by timestamp
func (r *conversationRepository) GetPreviousConversation(ctx context.Context, currentID uuid.UUID) (*models.Conversation, error) {
	query := `
		SELECT id, title, is_favorite, created_at
		FROM conversations
		WHERE created_at < (SELECT created_at FROM conversations WHERE id = $1)
		ORDER BY created_at DESC
		LIMIT 1
	`

	var conversation models.Conversation
	err := r.db.QueryRowContext(ctx, query, currentID).Scan(
		&conversation.ID,
		&conversation.Title,
		&conversation.IsFavorite,
		&conversation.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no previous conversation")
		}
		return nil, fmt.Errorf("failed to get previous conversation: %w", err)
	}

	return &conversation, nil
}
