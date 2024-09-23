package simworks

import (
	"testing"
)

func TestGenerateRandomID(t *testing.T) {
	id := GenerateRandomID(12)
	if len(id) != 12 {
		t.Errorf("Expected ID length to be 12, got %d", len(id))
	}

	id = GenerateRandomID(100)
	if len(id) != 100 {
		t.Errorf("Expected ID length to be 100, got %d", len(id))
	}

	id = GenerateRandomID(0)
	if len(id) != 0 {
		t.Errorf("Expected ID length to be 0, got %d", len(id))
	}

	id = GenerateRandomID(-42)
	if id != "" {
		t.Errorf("Expected ID to be an empty string for negative length, got %s", id)
	}
}
