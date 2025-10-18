package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/ventrosky/prompt-explorer/business/models"
)

// Repository defines the interface for all data operations
type Repository interface {
	// Conversation operations
	CreateConversation(ctx context.Context, conversation *models.Conversation) error
	GetConversationByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error)
	ListConversations(ctx context.Context) ([]*models.Conversation, error)
	UpdateConversation(ctx context.Context, conversation *models.Conversation) error
	DeleteConversation(ctx context.Context, id uuid.UUID) error
	ToggleConversationFavorite(ctx context.Context, id uuid.UUID) error
	GetNextConversation(ctx context.Context, currentID uuid.UUID) (*models.Conversation, error)
	GetPreviousConversation(ctx context.Context, currentID uuid.UUID) (*models.Conversation, error)

	// Message operations
	CreateMessage(ctx context.Context, message *models.Message) error
	GetMessageByID(ctx context.Context, id uuid.UUID) (*models.Message, error)
	GetMessagesByConversationID(ctx context.Context, conversationID uuid.UUID) ([]*models.Message, error)
	GetSystemPrompt(ctx context.Context, conversationID uuid.UUID) (string, error)
	UpdateMessage(ctx context.Context, message *models.Message) error
	DeleteMessage(ctx context.Context, id uuid.UUID) error
	DeleteMessagesByConversationID(ctx context.Context, conversationID uuid.UUID) error
}

// repository implements the Repository interface
type repository struct {
	db *sql.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// Conversation operations

func (r *repository) CreateConversation(ctx context.Context, conversation *models.Conversation) error {
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

func (r *repository) GetConversationByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error) {
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

func (r *repository) ListConversations(ctx context.Context) ([]*models.Conversation, error) {
	query := `
		SELECT id, title, is_favorite, created_at
		FROM conversations
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query conversations: %w", err)
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

func (r *repository) UpdateConversation(ctx context.Context, conversation *models.Conversation) error {
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

func (r *repository) DeleteConversation(ctx context.Context, id uuid.UUID) error {
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

func (r *repository) ToggleConversationFavorite(ctx context.Context, id uuid.UUID) error {
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

func (r *repository) GetNextConversation(ctx context.Context, currentID uuid.UUID) (*models.Conversation, error) {
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

func (r *repository) GetPreviousConversation(ctx context.Context, currentID uuid.UUID) (*models.Conversation, error) {
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

// Message operations

func (r *repository) CreateMessage(ctx context.Context, message *models.Message) error {
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

func (r *repository) GetMessageByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
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

func (r *repository) GetMessagesByConversationID(ctx context.Context, conversationID uuid.UUID) ([]*models.Message, error) {
	query := `
		SELECT id, conversation_id, sender, content, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
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

func (r *repository) GetSystemPrompt(ctx context.Context, conversationID uuid.UUID) (string, error) {
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

func (r *repository) UpdateMessage(ctx context.Context, message *models.Message) error {
	query := `
		UPDATE messages
		SET content = $2
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		message.ID,
		message.Content,
	)

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

func (r *repository) DeleteMessage(ctx context.Context, id uuid.UUID) error {
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

func (r *repository) DeleteMessagesByConversationID(ctx context.Context, conversationID uuid.UUID) error {
	query := `DELETE FROM messages WHERE conversation_id = $1`

	_, err := r.db.ExecContext(ctx, query, conversationID)
	if err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}

	return nil
}


