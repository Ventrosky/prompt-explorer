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
