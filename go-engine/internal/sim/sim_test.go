package sim

import (
	"testing"
	"time"
)

func TestNewSimulation(t *testing.T) {
	name := "Test Simulation"
	motto := "Safety First"
	sim := NewSimulation(name, motto)

	if sim == nil {
		t.Fatal("NewSimulation returned nil")
	}

	// Check SimInfo
	if sim.info.Name != name {
		t.Errorf("Expected name %s, got %s", name, sim.info.Name)
	}
	if sim.info.Motto != motto {
		t.Errorf("Expected motto %s, got %s", motto, sim.info.Motto)
	}
	if sim.info.ID == "" {
		t.Error("Expected non-empty ID")
	}
	if sim.info.SpawnedAt.IsZero() {
		t.Error("Expected non-zero SpawnedAt time")
	}

	// Check default clock settings
	expectedStartTime := time.Date(2000, 1, 1, 8, 0, 0, 0, time.FixedZone("EST", -5*60*60))
	if !sim.clock.startedAt.Equal(expectedStartTime) {
		t.Errorf("Expected start time %v, got %v", expectedStartTime, sim.clock.startedAt)
	}
	if sim.clock.currentIter != 0 {
		t.Errorf("Expected currentIter 0, got %d", sim.clock.currentIter)
	}

	// Check default environment
	if sim.environment.Weather != "Sunny" {
		t.Errorf("Expected default weather 'Sunny', got %s", sim.environment.Weather)
	}

	// Check other default settings
	if len(sim.components) != 0 {
		t.Errorf("Expected 0 components, got %d", len(sim.components))
	}
	if sim.running {
		t.Error("Expected simulation to not be running")
	}
}

func TestSimulationAdvanceAndStop(t *testing.T) {
	// Initialize a simulation
	sim := NewSimulation("Test Sim", "Safety First")

	// Create a boiler and start it
	boiler := NewBoiler("Main Boiler")
	boiler.TurnOn()

	// Add the boiler
	sim.AddComponent(boiler)

	// Create and add a turbine
	turbine := NewTurbine("Main Turbine")
	sim.AddComponent(turbine)

	// Advance the simulation 1000000 steps
	sim.Advance(1000000)

	// Wait a short time to ensure the simulation has started
	time.Sleep(2 * time.Millisecond)

	// Call Stop
	sim.Stop()

	// Wait for the simulation to fully stop
	for sim.IsRunning() {
		time.Sleep(10 * time.Millisecond)
	}

	// Check that the number of iterations is less than 1000000
	if sim.clock.currentIter >= 1000000 {
		t.Errorf("Expected iterations to be less than 1000000, got %d", sim.clock.currentIter)
	}

	// Additional checks to ensure the simulation ran and stopped correctly
	if sim.clock.currentIter == 0 {
		t.Error("Simulation did not advance at all")
	}

	if sim.IsRunning() {
		t.Error("Simulation should not be running after stop")
	}
}
