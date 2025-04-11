package telegram

import (
	"context"
	"fmt"
)

// MockTelegramClient is a mock implementation of TelegramClient for testing
type MockTelegramClient struct {
	// Add fields to control mock behavior if needed, e.g., expected errors, return values
}

// NewMockTelegramClient creates a new instance of MockTelegramClient
func NewMockTelegramClient() *MockTelegramClient {
	return &MockTelegramClient{}
}

// Authorize implements the TelegramClient interface for the mock
func (m *MockTelegramClient) Authorize(ctx context.Context) error {
	fmt.Println("MockTelegramClient: Authorize called")
	// Add mock logic here, e.g., return m.ExpectedAuthorizeError
	return nil
}

// FetchMessages implements the TelegramClient interface for the mock
func (m *MockTelegramClient) FetchMessages(ctx context.Context, chatID int64, from, to int64) ([]Message, error) {
	fmt.Printf("MockTelegramClient: FetchMessages called for chat %d\n", chatID)
	// Add mock logic here, e.g., return m.ExpectedMessages, m.ExpectedFetchError
	return []Message{}, nil
}

// ListChats implements the TelegramClient interface for the mock
func (m *MockTelegramClient) ListChats(ctx context.Context) ([]int64, error) {
	fmt.Println("MockTelegramClient: ListChats called")
	// Add mock logic here, e.g., return m.ExpectedChatIDs, m.ExpectedListChatsError
	return []int64{}, nil
}

// Example test using the mock (optional, for demonstration)
/*
func TestSomethingUsingTelegramClient(t *testing.T) {
	mockClient := NewMockTelegramClient()

	// Inject mockClient into the component that uses TelegramClient
	// component := NewComponent(mockClient)

	// Call the component's method that uses the client
	// err := component.DoSomething(context.Background())

	// Assert expectations, e.g., check if mock methods were called, check results
	// if err != nil {
	// 	t.Errorf("Expected no error, got %v", err)
	// }
}
*/
