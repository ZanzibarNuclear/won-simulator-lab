package simworks

import (
	"testing"
)

func TestSimWorks(t *testing.T) {
	s := NewSimulator("Test Simulator", "Test Simulator", "Test Simulator")

	if s.ID != "Test Simulator" {
		t.Errorf("Expected ID to be 'Test Simulator', got %s", s.ID)
	}

	if s.Name != "Test Simulator" {
		t.Errorf("Expected Name to be 'Test Simulator', got %s", s.Name)
	}

	if s.Purpose != "Test Simulator" {
		t.Errorf("Expected Purpose to be 'Test Simulator', got %s", s.Purpose)
	}

	if s.CreationDate.IsZero() {
		t.Errorf("Expected CreationDate to be non-zero, got %s", s.CreationDate)
	}

	if s.Clock == nil {
		t.Errorf("Expected Clock to be non-nil, got %v", s.Clock)
	}

	if s.Environment == nil {
		t.Errorf("Expected Environment to be non-nil, got %v", s.Environment)
	}

}
