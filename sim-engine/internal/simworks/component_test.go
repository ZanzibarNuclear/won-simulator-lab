package simworks

import (
	"testing"
)

func TestBaseComponent(t *testing.T) {
	bc := NewBaseComponent("Test Component", "Test Description")

	if bc.ID() == "" {
		t.Errorf("Expected ID to be non-empty, got %s", bc.ID())
	}

	if bc.Name() != "Test Component" {
		t.Errorf("Expected name to be 'Test Component', got %s", bc.Name())
	}

	if bc.Description() != "Test Description" {
		t.Errorf("Expected description to be 'Test Description', got %s", bc.Description())
	}

	status := bc.Status()

	if status["ID"] == nil {
		t.Errorf("Expected status to contain ID, got %v", status["ID"])
	}

	if status["Name"] != "Test Component" {
		t.Errorf("Expected status to contain Name, got %v", status["Name"])
	}

	if status["Description"] != "Test Description" {
		t.Errorf("Expected status to contain Description, got %v", status["Description"])
	}

}

func TestBaseComponentUpdateChangesLatestMoment(t *testing.T) {
	bc := NewBaseComponent("Test Component", "Test Description")
	sim := NewSimulator("Test Simulator", "Test Simulator", "Test Simulator")
	sim.AddComponent(bc)
	sim.Step()

	if bc.LatestMoment().IsZero() {
		t.Errorf("Expected LatestMoment to be non-zero, got %v", bc.LatestMoment())
	}

	sim.Step()
}
