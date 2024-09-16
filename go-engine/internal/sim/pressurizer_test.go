package sim

import (
	"math"
	"testing"
)

func TestPressurizer(t *testing.T) {
	// Create a new pressurizer
	pressurizer := NewPressurizer("TestPressurizerInitialConditions")

	// Create a new simulation and environment
	sim := NewSimulation("TestPressurizer", "TestEnvironment")
	env := NewEnvironment()

	// Add the pressurizer to the simulation
	sim.AddComponent(pressurizer)

	// Perform basic checks
	if pressurizer.GetName() != "TestPressurizerInitialConditions" {
		t.Errorf("Expected pressurizer name to be 'TestPressurizerInitialConditions', got '%s'", pressurizer.GetName())
	}

	if pressurizer.pressure != 0.0 {
		t.Errorf("Expected initial pressure to be 0.0, got %f", pressurizer.pressure)
	}

	if pressurizer.temperature != ROOM_TEMPERATURE {
		t.Errorf("Expected initial temperature to be %f, got %f", ROOM_TEMPERATURE, pressurizer.temperature)
	}

	// Test update function
	pressurizer.Update(env, sim)

	// Add more specific tests here based on the expected behavior of the Pressurizer
}

func TestPressurizerReachesTargetPT(t *testing.T) {
	// Create a new simulation, environment, and pressurizer
	sim := NewSimulation("TestPressurizer", "TestEnvironment")
	env := NewEnvironment()
	pressurizer := NewPressurizer("TestPressurizerReachesTargetPT")

	// Add the pressurizer to the simulation
	sim.AddComponent(pressurizer)

	// Turn on the heater
	pressurizer.heaterOn = true

	// Update 100 times
	for i := 0; i < 10; i++ {
		pressurizer.Update(env, sim)
		pressurizer.PrintStatus()
	}

	// Check that pressure is at target pressure
	if !almostEqual(pressurizer.pressure, pressurizer.targetPressure, 0.1) {
		t.Errorf("Expected pressure to be %f, got %f", pressurizer.targetPressure, pressurizer.pressure)
	}

	// Check that temperature is at target temperature
	if !almostEqual(pressurizer.temperature, TARGET_TEMPERATURE, 0.1) {
		t.Errorf("Expected temperature to be %f, got %f", TARGET_TEMPERATURE, pressurizer.temperature)
	}
}

// Helper function to compare float values with a tolerance
func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
