package pwr

import (
	"testing"

	"worldofnuclear.com/internal/simworks"
)

func TestNewPrimaryLoop(t *testing.T) {
	pl := NewPrimaryLoop("TestLoop-Pump", "The is a test.")
	if pl.Name() != "TestLoop-Pump" {
		t.Errorf("Expected name to be TestLoop-Pump, got %s", pl.Name())
	}
	if pl.Description() != "The is a test." {
		t.Errorf("Expected description to be The is a test., got %s", pl.Description())
	}
}

func TestSimulatorDrivesPrimaryLoop(t *testing.T) {
	pl := NewPrimaryLoop("TestLoop-Pump", "The is a test.")
	sim := simworks.NewSimulator("Test Sim", "Test Sim")
	sim.AddComponent(pl)
	sim.Run(100)

	if pl.LatestMoment().IsZero() {
		t.Errorf("Expected LatestMoment to be non-zero, got %v", pl.LatestMoment())
	}
}
