package summarizer

import (
	"fmt"
	"testing"

	// Adjust the import path according to your module name
	"github.com/azalio/tg-summary/internal/telegram"
)

// MockSummarizer is a mock implementation of Summarizer for testing
type MockSummarizer struct {
	// Add fields to control mock behavior if needed
	ExpectedSummary string
	ExpectedError   error
}

// NewMockSummarizer creates a new instance of MockSummarizer
func NewMockSummarizer() *MockSummarizer {
	return &MockSummarizer{
		ExpectedSummary: "This is a mock summary.",
	}
}

// Summarize implements the Summarizer interface for the mock
func (m *MockSummarizer) Summarize(messages []telegram.Message) (string, error) {
	fmt.Printf("MockSummarizer: Summarize called with %d messages\n", len(messages))
	if m.ExpectedError != nil {
		return "", m.ExpectedError
	}
	return m.ExpectedSummary, nil
}

// Example test using the mock (keep testing import)
func TestSummarizerMock(t *testing.T) {
	mockSummarizer := NewMockSummarizer()
	summary, err := mockSummarizer.Summarize([]telegram.Message{{ID: 1}, {ID: 2}})
	if err != nil {
		t.Errorf("Summarize failed: %v", err)
	}
	if summary != mockSummarizer.ExpectedSummary {
		t.Errorf("Expected summary '%s', got '%s'", mockSummarizer.ExpectedSummary, summary)
	}
}