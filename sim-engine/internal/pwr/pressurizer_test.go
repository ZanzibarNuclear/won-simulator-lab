package pwr

import (
	"testing"

	"worldofnuclear.com/internal/simworks"
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

func TestHeaterPowerEvent(t *testing.T) {
	p := NewPressurizer("TestPressurizer", "Test pressurizer description")
	pwrSim := NewPwrSim("Test PWR", "Heater power test simulation")
	pwrSim.AddComponent(p)

	// Initially, the heater should be off
	if p.HeaterOn() {
		t.Error("Expected heater to be off initially")
	}

	// Create and queue the event to turn the heater on
	eventToWatch := NewEvent_HeaterPower(true)
	pwrSim.QueueEvent(eventToWatch)
	pwrSim.Step()

	if !p.HeaterOn() {
		t.Error("Expected heater to be on after the event")
	}
	if !eventToWatch.IsComplete() {
		t.Error("HeaterPowerEvent should be completed after one step")
	}

	// Now test turning the heater off
	eventToWatch = NewEvent_HeaterPower(false)
	pwrSim.QueueEvent(eventToWatch)
	pwrSim.Step()

	if p.HeaterOn() {
		t.Error("Expected heater to be off after the second event")
	}
	if !eventToWatch.IsComplete() {
		t.Error("HeaterPowerEvent (turn off) should be completed after one step")
	}
}

func TestSprayNozzleEvent(t *testing.T) {
	p := NewPressurizer("TestPressurizer", "Test pressurizer description")
	pwrSim := NewPwrSim("Test PWR", "Spray nozzle test simulation")
	pwrSim.AddComponent(p)

	// Initially, the spray nozzle should be closed
	if p.SprayNozzleOpen() {
		t.Error("Expected spray nozzle to be closed initially")
	}

	// Create and queue the event to open the spray nozzle
	eventToWatch := NewEvent_SprayNozzle(true)
	pwrSim.QueueEvent(eventToWatch)
	pwrSim.Step()

	if !p.SprayNozzleOpen() {
		t.Error("Expected spray nozzle to be open after the event")
	}
	if p.SprayFlowRate() == 0.0 {
		t.Error("Expected spray flow rate to be non-zero when nozzle is open")
	}
	if !eventToWatch.IsComplete() {
		t.Error("SprayNozzleEvent should be completed after one step")
	}

	// Check if spray flow rate is set correctly
	expectedSprayFlowRate := Config["pressurizer"]["spray_flow_rate"]
	if p.SprayFlowRate() != expectedSprayFlowRate {
		t.Errorf("Expected spray flow rate to be %f, got %f", expectedSprayFlowRate, p.SprayFlowRate())
	}

	// Now test closing the spray nozzle
	eventToWatch = NewEvent_SprayNozzle(false)
	pwrSim.QueueEvent(eventToWatch)
	pwrSim.Step()

	if p.SprayNozzleOpen() {
		t.Error("Expected spray nozzle to be closed after the second event")
	}
	if p.SprayFlowRate() != 0.0 {
		t.Errorf("Expected spray flow rate to be 0.0 when nozzle is closed, got %f", p.SprayFlowRate())
	}
	if !eventToWatch.IsComplete() {
		t.Error("SprayNozzleEvent (close) should be completed after one step")
	}

	// Check if spray flow rate is zero when nozzle is closed
	if p.SprayFlowRate() != 0.0 {
		t.Errorf("Expected spray flow rate to be 0.0 when nozzle is closed, got %f", p.SprayFlowRate())
	}
}

func TestSetPressureEvent(t *testing.T) {
	p := NewPressurizer("TestPressurizer", "Test pressurizer description")
	pwrSim := NewPwrSim("Test PWR", "We test so you don't have to.")
	pwrSim.AddComponent(p)

	initialPressure := p.Pressure()
	targetPressure := 15.5 // MPa, typical operating pressure

	// Create and queue the event
	eventToWatch := NewEvent_TargetPressure(targetPressure)
	pwrSim.QueueEvent(eventToWatch)
	pwrSim.QueueEvent(NewEvent_HeaterPower(true))

	pwrSim.RunForABit(0, 0, 0, 10)
	if eventToWatch.IsInProgress() && !p.HeaterOnHigh() {
		t.Error("Heater should be on high power to raise pressure")
	}

	pwrSim.RunForABit(1, 0, 0, 0)

	if p.Pressure() == initialPressure {
		t.Errorf("Pressure did not change. Expected it to move towards %f, but got %f", targetPressure, p.Pressure())
	}
	if p.Pressure() < targetPressure {
		t.Errorf("Pressure did not reach target. Expected %f, got %f", targetPressure, p.Pressure())
	}
	if !eventToWatch.IsComplete() {
		t.Error("Event should be completed after reaching target pressure")
	}
	if eventToWatch.IsComplete() && !p.HeaterOnLow() {
		t.Error("Heater should be on low power after reaching target pressure")
	}
}

func TestReliefValveVentEvent(t *testing.T) {
	p := NewPressurizer("TestPressurizer", "Test pressurizer description")
	pwrSim := NewPwrSim("Test PWR", "Relief valve vent test simulation")
	pwrSim.AddComponent(p)
	pwrSim.SetEventHandler(pwrSim)

	targetPressure := 20.0 // higher than threshold

	// Set up conditions to raise pressure
	pressureEvent := NewEvent_TargetPressure(targetPressure)
	pwrSim.QueueEvent(pressureEvent)
	heaterEvent := NewEvent_HeaterPower(true)
	pwrSim.QueueEvent(heaterEvent)
	pwrSim.Step()

	if !pressureEvent.IsInProgress() {
		t.Error("Pressure event should be in progress")
	}
	if !heaterEvent.IsComplete() {
		t.Error("Heater event should be complete")
	}

	pwrSim.RunForABit(0, 0, 0, 5)

	if len(pwrSim.Events) > 1 {
		t.Errorf("Expected 1 event at this point, got %d", len(pwrSim.Events))
		return
	}
	if len(pwrSim.InactiveEvents) != 1 {
		t.Errorf("Expected 1 inactive event at this point, got %d", len(pwrSim.InactiveEvents))
		return
	}

	pwrSim.RunForABit(0, 0, 1, 2)
	// t.Logf("Pressurizer status: %v", p.Status())
	// for _, event := range pwrSim.Events {
	// 	t.Logf("Event: %v", event)
	// }
	// for _, event := range pwrSim.InactiveEvents {
	// 	t.Logf("Inactive event: %v", event)
	// }

	if len(pwrSim.InactiveEvents) != 2 {
		t.Errorf("Expected 2 inactive event at this point, got %d", len(pwrSim.InactiveEvents))
		return
	}

	// Find the relief valve event in InactiveEvents
	var reliefValveEvent *simworks.Event
	for _, event := range pwrSim.InactiveEvents {
		// t.Logf("Inactive event: %v", event)
		if event.Code == Event_pr_reliefValveVent {
			reliefValveEvent = event
			break
		}
	}

	if reliefValveEvent == nil {
		t.Error("Relief valve vent event not found in InactiveEvents")
		return
	}
	if !reliefValveEvent.IsComplete() {
		t.Error("Relief valve vent event should have been processed by simulator")
	}
}
