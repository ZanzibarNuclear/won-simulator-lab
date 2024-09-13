package sim

import (
	"fmt"
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

	// Add one of each component type
	sim.AddComponent(NewPrimaryLoop("Lara's Primary Loop"))
	sim.AddComponent(NewSecondaryLoop("Lana's Secondary Loop"))
	sim.AddComponent(NewReactorCore("Corey the Reactor Core"))
	sim.AddComponent(NewPressurizer("#1 Pressurizer"))
	sim.AddComponent(NewSteamGenerator("Stevie the Steam Generator"))
	sim.AddComponent(NewSteamTurbine("Tilly the Steam Turbine"))
	sim.AddComponent(NewCondenser("Condie the Condenser"))
	sim.AddComponent(NewGenerator("Flo the Power Generator"))

	// turn things on so that update loop is meaningful
	sim.FindPrimaryLoop().SwitchOnPump()
	// TODO: turn on more as component models improve

	// Verify that components were added
	expectedComponents := 8
	if len(sim.components) != expectedComponents {
		t.Errorf("Expected %d components, got %d", expectedComponents, len(sim.components))
	}

	// Verify each component type
	componentTypes := map[string]bool{
		"*sim.PrimaryLoop":    false,
		"*sim.SecondaryLoop":  false,
		"*sim.ReactorCore":    false,
		"*sim.Pressurizer":    false,
		"*sim.SteamGenerator": false,
		"*sim.SteamTurbine":   false,
		"*sim.Generator":      false,
		"*sim.Condenser":      false,
	}

	for _, component := range sim.components {
		componentType := fmt.Sprintf("%T", component)
		if _, exists := componentTypes[componentType]; exists {
			componentTypes[componentType] = true
		}
	}

	for componentType, found := range componentTypes {
		if !found {
			t.Errorf("Expected component of type %s, but it was not found", componentType)
		}
	}
	// Advance the simulation 1000000 steps
	sim.Advance(1000000)
	time.Sleep(2 * time.Millisecond)
	sim.Stop()

	time.Sleep(2 * time.Millisecond)
	if sim.IsRunning() {
		t.Error("Simulation should not be running after stop")
	}

	// Check that the number of iterations is less than 1000000
	if sim.clock.currentIter >= 1000000 {
		t.Errorf("Expected interruption before completion, got through all %d ticks", sim.clock.currentIter)
	}

	// Additional checks to ensure the simulation ran and stopped correctly
	if sim.clock.currentIter == 0 {
		t.Error("Simulation did not advance at all")
	}
}
