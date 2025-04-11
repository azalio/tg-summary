package delivery

import (
	"fmt"
	"testing"
)

// MockDigestSender is a mock implementation of DigestSender for testing
type MockDigestSender struct {
	// Add fields to control mock behavior if needed
	SentDigests map[int64]string
	SendError   error
}

// NewMockDigestSender creates a new instance of MockDigestSender
func NewMockDigestSender() *MockDigestSender {
	return &MockDigestSender{
		SentDigests: make(map[int64]string),
	}
}

// SendDigest implements the DigestSender interface for the mock
func (m *MockDigestSender) SendDigest(chatID int64, digest string) error {
	fmt.Printf("MockDigestSender: SendDigest called for chat %d\n", chatID)
	if m.SendError != nil {
		return m.SendError
	}
	m.SentDigests[chatID] = digest
	return nil
}

// Example test using the mock (keep testing import)
func TestDeliveryMock(t *testing.T) {
	mockSender := NewMockDigestSender()
	err := mockSender.SendDigest(123, "Test digest")
	if err != nil {
		t.Errorf("SendDigest failed: %v", err)
	}
	if _, ok := mockSender.SentDigests[123]; !ok {
		t.Errorf("Digest for chat 123 was not sent")
	}
}