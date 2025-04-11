package scheduler

import (
	"fmt"
	"testing"
)

// MockScheduler is a mock implementation of Scheduler for testing
type MockScheduler struct {
	// Add fields to control mock behavior if needed
	Started bool
	Stopped bool
	StartError error
	StopError  error
}

// NewMockScheduler creates a new instance of MockScheduler
func NewMockScheduler() *MockScheduler {
	return &MockScheduler{}
}

// Start implements the Scheduler interface for the mock
func (m *MockScheduler) Start() error {
	fmt.Println("MockScheduler: Start called")
	if m.StartError != nil {
		return m.StartError
	}
	m.Started = true
	m.Stopped = false
	return nil
}

// Stop implements the Scheduler interface for the mock
func (m *MockScheduler) Stop() error {
	fmt.Println("MockScheduler: Stop called")
	if m.StopError != nil {
		return m.StopError
	}
	m.Started = false
	m.Stopped = true
	return nil
}

// Example test using the mock (keep testing import)
func TestSchedulerMock(t *testing.T) {
	mockScheduler := NewMockScheduler()
	err := mockScheduler.Start()
	if err != nil {
		t.Errorf("Start failed: %v", err)
	}
	if !mockScheduler.Started {
		t.Errorf("Scheduler was not started")
	}
	err = mockScheduler.Stop()
	if err != nil {
		t.Errorf("Stop failed: %v", err)
	}
	if !mockScheduler.Stopped {
		t.Errorf("Scheduler was not stopped")
	}
}