package simworks

import (
	"testing"
	"time"
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

	latestMoment, ok := status["LatestMoment"].(time.Time)
	if !ok {
		t.Error("Expected to find LatestMoment in status, but it was not found")
	}
	if !latestMoment.IsZero() {
		t.Errorf("Expected LatestMoment to be zero before first update, got %v", latestMoment)
	}
}

func TestBaseComponentUpdateChangesLatestMoment(t *testing.T) {
	bc := NewBaseComponent("Test Component", "Test Description")
	sim := NewSimulator("Test Simulator", "Test Simulator", "Test Simulator")
	sim.AddComponent(bc)
	sim.Step()

	marker := bc.LatestMoment()
	simNow := sim.Clock.SimNow()

	if marker.IsZero() {
		t.Errorf("Expected LatestMoment to be non-zero, got %v", bc.LatestMoment())
	}

	if simNow != marker {
		t.Errorf("SimNow should be in sync with LatestMoment, expected %v, got %v", marker, simNow)
	}

	sim.Run(90)
	delta := bc.LatestMoment().Sub(marker)

	if delta.Seconds() != 90 {
		t.Errorf("Expected LatestMoment to be 90 seconds after first update, got %v", delta)
	}

}
