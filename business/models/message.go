package models

import (
	"time"

	"github.com/google/uuid"
)

// MessageType represents the type of message sender
type MessageType string

const (
	MessageTypeSystem    MessageType = "system"
	MessageTypeUser      MessageType = "user"
	MessageTypeAssistant MessageType = "assistant"
)

// Message represents a message in a conversation
type Message struct {
	ID             uuid.UUID   `json:"id" db:"id"`
	ConversationID uuid.UUID   `json:"conversation_id" db:"conversation_id"`
	Sender         MessageType `json:"sender" db:"sender"`
	Content        string      `json:"content" db:"content"`
	CreatedAt      time.Time   `json:"created_at" db:"created_at"`
}

// NewMessage creates a new message
func NewMessage(conversationID uuid.UUID, sender MessageType, content string) *Message {
	return &Message{
		ID:             uuid.New(),
		ConversationID: conversationID,
		Sender:         sender,
		Content:        content,
		CreatedAt:      time.Now(),
	}
}
