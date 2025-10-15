package services

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ventrosky/prompt-explorer/business/models"
	"github.com/ventrosky/prompt-explorer/business/repositories"
)

// TestDB represents a test database connection
type TestDB struct {
	DB                  *sql.DB
	ConversationRepo    repositories.ConversationRepository
	MessageRepo         repositories.MessageRepository
	ConversationService *ConversationService
}

// setupTestDB sets up a test database connection
func setupTestDB(t *testing.T) *TestDB {

	// Connect to test database
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/test_prompt_explorer?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	// Create repositories
	conversationRepo := repositories.NewConversationRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	// Create services
	conversationService := NewConversationService(conversationRepo, messageRepo)

	return &TestDB{
		DB:                  db,
		ConversationRepo:    conversationRepo,
		MessageRepo:         messageRepo,
		ConversationService: conversationService,
	}
}

// cleanupTestDB cleans up test database
func (tdb *TestDB) cleanupTestDB(t *testing.T) {
	if tdb.DB != nil {
		// Clean up test data
		_, err := tdb.DB.Exec("DELETE FROM messages")
		if err != nil {
			t.Logf("Failed to clean messages: %v", err)
		}

		_, err = tdb.DB.Exec("DELETE FROM conversations")
		if err != nil {
			t.Logf("Failed to clean conversations: %v", err)
		}

		tdb.DB.Close()
	}
}

func TestConversationRepository_Create(t *testing.T) {
	tdb := setupTestDB(t)
	defer tdb.cleanupTestDB(t)

	ctx := context.Background()

	// Create a test conversation
	conversation := &models.Conversation{
		ID:         uuid.New(),
		Title:      "Test Conversation",
		IsFavorite: false,
		CreatedAt:  time.Now(),
	}

	// Test create
	err := tdb.ConversationRepo.Create(ctx, conversation)
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	// Verify conversation was created
	retrieved, err := tdb.ConversationRepo.GetByID(ctx, conversation.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve conversation: %v", err)
	}

	if retrieved.Title != conversation.Title {
		t.Errorf("Expected title %s, got %s", conversation.Title, retrieved.Title)
	}

	if retrieved.IsFavorite != conversation.IsFavorite {
		t.Errorf("Expected is_favorite %v, got %v", conversation.IsFavorite, retrieved.IsFavorite)
	}
}

func TestConversationRepository_List(t *testing.T) {
	tdb := setupTestDB(t)
	defer tdb.cleanupTestDB(t)

	ctx := context.Background()

	// Create multiple test conversations
	conversations := []*models.Conversation{
		{
			ID:         uuid.New(),
			Title:      "First Conversation",
			IsFavorite: false,
			CreatedAt:  time.Now().Add(-2 * time.Hour),
		},
		{
			ID:         uuid.New(),
			Title:      "Second Conversation",
			IsFavorite: true,
			CreatedAt:  time.Now().Add(-1 * time.Hour),
		},
		{
			ID:         uuid.New(),
			Title:      "Third Conversation",
			IsFavorite: false,
			CreatedAt:  time.Now(),
		},
	}

	// Create conversations
	for _, conv := range conversations {
		err := tdb.ConversationRepo.Create(ctx, conv)
		if err != nil {
			t.Fatalf("Failed to create conversation: %v", err)
		}
	}

	// Test list
	list, err := tdb.ConversationRepo.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list conversations: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("Expected 3 conversations, got %d", len(list))
	}

	// Verify ordering (should be newest first)
	if list[0].Title != "Third Conversation" {
		t.Errorf("Expected first conversation to be 'Third Conversation', got %s", list[0].Title)
	}
}

func TestConversationRepository_ToggleFavorite(t *testing.T) {
	tdb := setupTestDB(t)
	defer tdb.cleanupTestDB(t)

	ctx := context.Background()

	// Create a test conversation
	conversation := &models.Conversation{
		ID:         uuid.New(),
		Title:      "Test Conversation",
		IsFavorite: false,
		CreatedAt:  time.Now(),
	}

	err := tdb.ConversationRepo.Create(ctx, conversation)
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	// Toggle favorite
	err = tdb.ConversationRepo.ToggleFavorite(ctx, conversation.ID)
	if err != nil {
		t.Fatalf("Failed to toggle favorite: %v", err)
	}

	// Verify favorite was toggled
	retrieved, err := tdb.ConversationRepo.GetByID(ctx, conversation.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve conversation: %v", err)
	}

	if !retrieved.IsFavorite {
		t.Error("Expected conversation to be favorited after toggle")
	}
}

func TestConversationService_CreateConversation(t *testing.T) {
	tdb := setupTestDB(t)
	defer tdb.cleanupTestDB(t)

	ctx := context.Background()

	// Test service create conversation
	conversation, err := tdb.ConversationService.CreateConversation(ctx, "Service Test Conversation")
	if err != nil {
		t.Fatalf("Failed to create conversation via service: %v", err)
	}

	if conversation.Title != "Service Test Conversation" {
		t.Errorf("Expected title 'Service Test Conversation', got %s", conversation.Title)
	}

	if conversation.IsFavorite {
		t.Error("Expected conversation to not be favorited by default")
	}
}

func TestConversationService_AddMessage(t *testing.T) {
	tdb := setupTestDB(t)
	defer tdb.cleanupTestDB(t)

	ctx := context.Background()

	// Create a conversation first
	conversation, err := tdb.ConversationService.CreateConversation(ctx, "Test Conversation")
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	tests := []struct {
		name           string
		conversationID string
		sender         models.MessageType
		content        string
		expectError    bool
		expectedSender models.MessageType
	}{
		{
			name:           "Add system message",
			conversationID: conversation.ID.String(),
			sender:         models.MessageTypeSystem,
			content:        "You are a helpful assistant",
			expectError:    false,
			expectedSender: models.MessageTypeSystem,
		},
		{
			name:           "Add user message",
			conversationID: conversation.ID.String(),
			sender:         models.MessageTypeUser,
			content:        "Hello, how are you?",
			expectError:    false,
			expectedSender: models.MessageTypeUser,
		},
		{
			name:           "Add assistant message",
			conversationID: conversation.ID.String(),
			sender:         models.MessageTypeAssistant,
			content:        "I'm doing well, thank you!",
			expectError:    false,
			expectedSender: models.MessageTypeAssistant,
		},
		{
			name:           "Add message to non-existent conversation",
			conversationID: "00000000-0000-0000-0000-000000000000",
			sender:         models.MessageTypeUser,
			content:        "Test",
			expectError:    true,
			expectedSender: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := tdb.ConversationService.AddMessage(ctx, tt.conversationID, tt.sender, tt.content)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if message.Sender != tt.expectedSender {
				t.Errorf("Expected sender %s, got %s", tt.expectedSender, message.Sender)
			}

			if message.Content != tt.content {
				t.Errorf("Expected content %s, got %s", tt.content, message.Content)
			}
		})
	}
}

func TestConversationService_UpdateMessage(t *testing.T) {
	tdb := setupTestDB(t)
	defer tdb.cleanupTestDB(t)

	ctx := context.Background()

	// Create conversation and message
	conversation, err := tdb.ConversationService.CreateConversation(ctx, "Test Conversation")
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	message, err := tdb.ConversationService.AddMessage(ctx, conversation.ID.String(), models.MessageTypeUser, "Original content")
	if err != nil {
		t.Fatalf("Failed to add message: %v", err)
	}

	tests := []struct {
		name        string
		messageID   string
		newContent  string
		expectError bool
	}{
		{
			name:        "Update existing message",
			messageID:   message.ID.String(),
			newContent:  "Updated content",
			expectError: false,
		},
		{
			name:        "Update non-existent message",
			messageID:   "00000000-0000-0000-0000-000000000000",
			newContent:  "Test",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tdb.ConversationService.UpdateMessage(ctx, tt.messageID, tt.newContent)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify the update
			updatedMessage, err := tdb.ConversationService.GetMessage(ctx, tt.messageID)
			if err != nil {
				t.Fatalf("Failed to get updated message: %v", err)
			}

			if updatedMessage.Content != tt.newContent {
				t.Errorf("Expected content %s, got %s", tt.newContent, updatedMessage.Content)
			}
		})
	}
}

func TestConversationService_DeleteMessage(t *testing.T) {
	tdb := setupTestDB(t)
	defer tdb.cleanupTestDB(t)

	ctx := context.Background()

	// Create conversation and message
	conversation, err := tdb.ConversationService.CreateConversation(ctx, "Test Conversation")
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	message, err := tdb.ConversationService.AddMessage(ctx, conversation.ID.String(), models.MessageTypeUser, "Test message")
	if err != nil {
		t.Fatalf("Failed to add message: %v", err)
	}

	tests := []struct {
		name        string
		messageID   string
		expectError bool
	}{
		{
			name:        "Delete existing message",
			messageID:   message.ID.String(),
			expectError: false,
		},
		{
			name:        "Delete non-existent message",
			messageID:   "00000000-0000-0000-0000-000000000000",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tdb.ConversationService.DeleteMessage(ctx, tt.messageID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify the message is deleted (only for successful deletions)
			if !tt.expectError {
				_, err = tdb.ConversationService.GetMessage(ctx, tt.messageID)
				if err == nil {
					t.Error("Expected error when getting deleted message")
				}
			}
		})
	}
}

func TestConversationService_GetMessage(t *testing.T) {
	tdb := setupTestDB(t)
	defer tdb.cleanupTestDB(t)

	ctx := context.Background()

	// Create conversation and message
	conversation, err := tdb.ConversationService.CreateConversation(ctx, "Test Conversation")
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	originalMessage, err := tdb.ConversationService.AddMessage(ctx, conversation.ID.String(), models.MessageTypeUser, "Test message content")
	if err != nil {
		t.Fatalf("Failed to add message: %v", err)
	}

	tests := []struct {
		name            string
		messageID       string
		expectError     bool
		expectedID      string
		expectedContent string
		expectedSender  models.MessageType
	}{
		{
			name:            "Get existing message",
			messageID:       originalMessage.ID.String(),
			expectError:     false,
			expectedID:      originalMessage.ID.String(),
			expectedContent: "Test message content",
			expectedSender:  models.MessageTypeUser,
		},
		{
			name:        "Get non-existent message",
			messageID:   "00000000-0000-0000-0000-000000000000",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message, err := tdb.ConversationService.GetMessage(ctx, tt.messageID)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if message.ID.String() != tt.expectedID {
				t.Errorf("Expected message ID %s, got %s", tt.expectedID, message.ID)
			}

			if message.Content != tt.expectedContent {
				t.Errorf("Expected content %s, got %s", tt.expectedContent, message.Content)
			}

			if message.Sender != tt.expectedSender {
				t.Errorf("Expected sender %s, got %s", tt.expectedSender, message.Sender)
			}
		})
	}
}

func TestConversationService_GetConversationHistory(t *testing.T) {
	tdb := setupTestDB(t)
	defer tdb.cleanupTestDB(t)

	ctx := context.Background()

	// Create conversation
	conversation, err := tdb.ConversationService.CreateConversation(ctx, "Test Conversation")
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}

	// Add multiple messages in order
	systemMsg, err := tdb.ConversationService.AddMessage(ctx, conversation.ID.String(), models.MessageTypeSystem, "System prompt")
	if err != nil {
		t.Fatalf("Failed to add system message: %v", err)
	}

	userMsg, err := tdb.ConversationService.AddMessage(ctx, conversation.ID.String(), models.MessageTypeUser, "User message")
	if err != nil {
		t.Fatalf("Failed to add user message: %v", err)
	}

	assistantMsg, err := tdb.ConversationService.AddMessage(ctx, conversation.ID.String(), models.MessageTypeAssistant, "Assistant response")
	if err != nil {
		t.Fatalf("Failed to add assistant message: %v", err)
	}

	// Get conversation history
	retrievedConversation, messages, err := tdb.ConversationService.GetConversationHistory(ctx, conversation.ID.String())
	if err != nil {
		t.Fatalf("Failed to get conversation history: %v", err)
	}

	if retrievedConversation.ID != conversation.ID {
		t.Errorf("Expected conversation ID %s, got %s", conversation.ID, retrievedConversation.ID)
	}

	if len(messages) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(messages))
	}

	// Verify message order (should be ordered by created_at ASC)
	if messages[0].ID != systemMsg.ID {
		t.Errorf("Expected first message to be system message")
	}

	if messages[1].ID != userMsg.ID {
		t.Errorf("Expected second message to be user message")
	}

	if messages[2].ID != assistantMsg.ID {
		t.Errorf("Expected third message to be assistant message")
	}

	// Verify message contents
	if messages[0].Content != "System prompt" {
		t.Errorf("Expected system message content 'System prompt', got %s", messages[0].Content)
	}

	if messages[1].Content != "User message" {
		t.Errorf("Expected user message content 'User message', got %s", messages[1].Content)
	}

	if messages[2].Content != "Assistant response" {
		t.Errorf("Expected assistant message content 'Assistant response', got %s", messages[2].Content)
	}
}
