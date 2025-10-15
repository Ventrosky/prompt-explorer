package models

import (
	"time"

	"github.com/google/uuid"
)

// Conversation represents a conversation in the system
type Conversation struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Title      string    `json:"title" db:"title"`
	IsFavorite bool      `json:"is_favorite" db:"is_favorite"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// NewConversation creates a new conversation
func NewConversation(title string) *Conversation {
	return &Conversation{
		ID:         uuid.New(),
		Title:      title,
		IsFavorite: false,
		CreatedAt:  time.Now(),
	}
}

// ToggleFavorite toggles the favorite status of the conversation
func (c *Conversation) ToggleFavorite() {
	c.IsFavorite = !c.IsFavorite
}
