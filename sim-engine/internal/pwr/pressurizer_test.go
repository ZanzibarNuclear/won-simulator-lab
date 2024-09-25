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
