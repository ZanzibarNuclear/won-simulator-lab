package simworks

import (
	"testing"
)

func TestBaseComponent(t *testing.T) {
	bc := NewBaseComponent("Test Component")

	if bc.ID == "" {
		t.Errorf("Expected ID to be non-empty, got %s", bc.ID)
	}

	if bc.Name != "Test Component" {
		t.Errorf("Expected name to be 'Test Component', got %s", bc.Name)
	}

	if _, ok := bc.GetStatus()["ID"]; !ok {
		t.Errorf("Expected status to contain ID, got %v", bc.GetStatus())
	}

	if _, ok := bc.GetStatus()["Name"]; !ok {
		t.Errorf("Expected status to contain Name, got %v", bc.GetStatus())
	}
}
