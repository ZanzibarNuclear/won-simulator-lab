package simworks

import (
	"testing"
)

func TestSimWorks(t *testing.T) {
	s := NewSimulator("Tess the Tester", "Test Simulator")

	if len(s.ID) != 8 {
		t.Errorf("Expected ID to be 8 characters, got %s", s.ID)
	}

	if s.Name != "Tess the Tester" {
		t.Errorf("Expected Name to be 'Tess the Tester', got %s", s.Name)
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

	if len(s.Components) != 0 {
		t.Errorf("Expected Components to be empty, got %v", s.Components)
	}

	if len(s.ComponentIndex) != 0 {
		t.Errorf("Expected ComponentIndex to be empty, got %v", s.ComponentIndex)
	}
}
