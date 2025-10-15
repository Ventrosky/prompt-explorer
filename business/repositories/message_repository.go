package repositories

import (
	"context"
	"database/sql"
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
	db *sql.DB
}

// NewMessageRepository creates a new MessageRepository
func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{
		db: db,
	}
}

// Create creates a new message
func (r *messageRepository) Create(ctx context.Context, message *models.Message) error {
	query := `
		INSERT INTO messages (id, conversation_id, sender, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		message.ID,
		message.ConversationID,
		message.Sender,
		message.Content,
		message.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

// GetByID retrieves a message by ID
func (r *messageRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	query := `
		SELECT id, conversation_id, sender, content, created_at
		FROM messages
		WHERE id = $1
	`

	var message models.Message
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&message.ID,
		&message.ConversationID,
		&message.Sender,
		&message.Content,
		&message.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return &message, nil
}

// GetByConversationID retrieves all messages for a conversation
func (r *messageRepository) GetByConversationID(ctx context.Context, conversationID uuid.UUID) ([]*models.Message, error) {
	query := `
		SELECT id, conversation_id, sender, content, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var message models.Message
		err := rows.Scan(
			&message.ID,
			&message.ConversationID,
			&message.Sender,
			&message.Content,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, &message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate messages: %w", err)
	}

	return messages, nil
}

// GetSystemPrompt retrieves the system prompt (first system message) for a conversation
func (r *messageRepository) GetSystemPrompt(ctx context.Context, conversationID uuid.UUID) (string, error) {
	query := `
		SELECT content
		FROM messages
		WHERE conversation_id = $1 AND sender = 'system'
		ORDER BY created_at ASC
		LIMIT 1
	`

	var content string
	err := r.db.QueryRowContext(ctx, query, conversationID).Scan(&content)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no system prompt found")
		}
		return "", fmt.Errorf("failed to get system prompt: %w", err)
	}

	return content, nil
}

// Update updates a message
func (r *messageRepository) Update(ctx context.Context, message *models.Message) error {
	query := `
		UPDATE messages
		SET content = $2, sender = $3
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, message.ID, message.Content, message.Sender)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

// Delete deletes a message
func (r *messageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM messages WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

// DeleteByConversationID deletes all messages for a conversation
func (r *messageRepository) DeleteByConversationID(ctx context.Context, conversationID uuid.UUID) error {
	query := `DELETE FROM messages WHERE conversation_id = $1`

	_, err := r.db.ExecContext(ctx, query, conversationID)
	if err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}

	return nil
}
