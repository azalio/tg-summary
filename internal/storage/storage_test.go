package storage

import (
	"fmt"
	"testing"

	// Adjust the import path according to your module name
	"github.com/azalio/tg-summary/internal/telegram" 
)

// MockMessageStorage is a mock implementation of MessageStorage for testing
type MockMessageStorage struct {
	// Add fields to control mock behavior if needed
	SavedMessages []telegram.Message
	GetError      error
	SaveError     error
}

// NewMockMessageStorage creates a new instance of MockMessageStorage
func NewMockMessageStorage() *MockMessageStorage {
	return &MockMessageStorage{
		SavedMessages: make([]telegram.Message, 0),
	}
}

// SaveMessage implements the MessageStorage interface for the mock
func (m *MockMessageStorage) SaveMessage(msg telegram.Message) error {
	fmt.Printf("MockMessageStorage: SaveMessage called for message ID %d\n", msg.ID)
	if m.SaveError != nil {
		return m.SaveError
	}
	m.SavedMessages = append(m.SavedMessages, msg)
	return nil
}

// GetMessages implements the MessageStorage interface for the mock
func (m *MockMessageStorage) GetMessages(chatID int64, from, to int64) ([]telegram.Message, error) {
	fmt.Printf("MockMessageStorage: GetMessages called for chat %d\n", chatID)
	if m.GetError != nil {
		return nil, m.GetError
	}
	// Basic filtering for mock, can be enhanced
	var result []telegram.Message
	for _, msg := range m.SavedMessages {
		if msg.ChatID == chatID && msg.Timestamp >= from && msg.Timestamp <= to {
			result = append(result, msg)
		}
	}
	return result, nil
}

// Example test using the mock (keep testing import)
func TestStorageMock(t *testing.T) {
	mockStorage := NewMockMessageStorage()
	// Example usage:
	err := mockStorage.SaveMessage(telegram.Message{ID: 1, ChatID: 123, Timestamp: 100})
	if err != nil {
		t.Errorf("SaveMessage failed: %v", err)
	}
	msgs, err := mockStorage.GetMessages(123, 50, 150)
	if err != nil {
		t.Errorf("GetMessages failed: %v", err)
	}
	if len(msgs) != 1 {
		t.Errorf("Expected 1 message, got %d", len(msgs))
	}
}