package pwr

import (
	"testing"
)

func TestNewPressurizerDefaults(t *testing.T) {
	p := NewPressurizer("TestPressurizer", "Test pressurizer description")

	// Test default values
	if len(p.ID()) == 0 {
		t.Errorf("Expected ID to be non-empty")
	}

	if p.Name() != "TestPressurizer" {
		t.Errorf("Expected name to be 'TestPressurizer', got %s", p.Name())
	}

	if p.Description() != "Test pressurizer description" {
		t.Errorf("Expected description to be 'Test pressurizer description', got %s", p.Description())
	}

	expectedPressure := Config["common"]["atmospheric_pressure"]
	if p.Pressure() != expectedPressure {
		t.Errorf("Expected default pressure to be %f, got %f", expectedPressure, p.Pressure())
	}

	expectedTemperature := Config["common"]["room_temperature"]
	if p.Temperature() != expectedTemperature {
		t.Errorf("Expected default temperature to be %f, got %f", expectedTemperature, p.Temperature())
	}

	if p.HeaterOn() {
		t.Error("Expected heater to be off by default")
	}

	if p.SprayNozzleOpen() {
		t.Error("Expected spray nozzle to be closed by default")
	}
}

func TestSetPressureEvent(t *testing.T) {
	p := NewPressurizer("TestPressurizer", "Test pressurizer description")
	pwrSim := NewPwrSim("Test PWR", "We test so you don't have to.")
	pwrSim.AddComponent(p)

	initialPressure := p.Pressure()
	targetPressure := 15.5 // MPa, typical operating pressure

	// Create and queue the event
	pwrSim.QueueEvent(NewTargetPressureEvent(targetPressure))
	eventToWatch := &pwrSim.Events[0]

	// Run the simulation for a few steps
	// Run the simulation for a few steps
	pwrSim.RunForABit(1, 0, 0, 0)

	// Check if the pressure has changed
	if p.Pressure() == initialPressure {
		t.Errorf("Pressure did not change. Expected it to move towards %f, but got %f", targetPressure, p.Pressure())
	}

	// Check if the pressure reached the target
	if p.Pressure() != targetPressure {
		t.Errorf("Pressure did not reach target. Expected %f, got %f", targetPressure, p.Pressure())
	}

	// Check if the event is completed
	if !eventToWatch.IsComplete() {
		t.Error("Event should be completed after reaching target pressure")
	}

	// Check if heater is on low power after reaching target
	if eventToWatch.IsComplete() && !p.HeaterOnLow() {
		t.Error("Heater should be on low power after reaching target pressure")
	}

	p.Print()
}
