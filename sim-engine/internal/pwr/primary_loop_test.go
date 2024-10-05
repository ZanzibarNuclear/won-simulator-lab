package pwr

import (
	"testing"

	"worldofnuclear.com/internal/simworks"
)

func TestPrimaryLoop_Init(t *testing.T) {
	pl := NewPrimaryLoop("TestLoop-Pump", "The is a test.")
	if pl.Name() != "TestLoop-Pump" {
		t.Errorf("Expected name to be TestLoop-Pump, got %s", pl.Name())
	}
	if pl.Description() != "The is a test." {
		t.Errorf("Expected description to be The is a test., got %s", pl.Description())
	}
}

func TestPrimaryLoop_SimulatorDrives(t *testing.T) {
	pl := NewPrimaryLoop("TestLoop-Pump", "The is a test.")
	sim := simworks.NewSimulator("Test Sim", "Test Sim")
	sim.AddComponent(pl)
	sim.Run(100)

	if pl.LatestMoment().IsZero() {
		t.Errorf("Expected LatestMoment to be non-zero, got %v", pl.LatestMoment())
	}
}

func TestPrimaryLoop_PumpSwitch(t *testing.T) {
	pl := NewPrimaryLoop("TestLoop-Pump", "The is a test.")
	sim := simworks.NewSimulator("Test Sim", "Test Sim")
	sim.AddComponent(pl)
	sim.QueueEvent(NewEvent_PumpSwitch(true))
	sim.Step()
	if !pl.PumpOn() {
		t.Errorf("Expected pump to be on, got %v", pl.PumpOn())
	}
	if pl.PumpPressure() != Config["primary_loop"]["pump_on_pressure"] {
		t.Errorf("Expected pump pressure to be %v, got %v", Config["primary_loop"]["pump_on_pressure"], pl.PumpPressure())
	}
	if pl.FlowRate() != Config["primary_loop"]["pump_on_flow_rate"] {
		t.Errorf("Expected pump flow rate to be %v, got %v", Config["primary_loop"]["pump_on_flow_rate"], pl.FlowRate())
	}
	if pl.PumpHeat() != Config["primary_loop"]["pump_on_heat"] {
		t.Errorf("Expected pump heat to be %v, got %v", Config["primary_loop"]["pump_on_heat"], pl.PumpHeat())
	}

	sim.QueueEvent(NewEvent_PumpSwitch(false))
	sim.Step()
	if pl.PumpOn() {
		t.Errorf("Expected pump to be off, got %v", pl.PumpOn())
	}
	if pl.PumpPressure() != Config["primary_loop"]["pump_off_pressure"] {
		t.Errorf("Expected pump pressure to be %v, got %v", Config["primary_loop"]["pump_off_pressure"], pl.PumpPressure())
	}
	if pl.FlowRate() != Config["primary_loop"]["pump_off_flow_rate"] {
		t.Errorf("Expected pump flow rate to be %v, got %v", Config["primary_loop"]["pump_off_flow_rate"], pl.FlowRate())
	}
	if pl.PumpHeat() != Config["primary_loop"]["pump_off_heat"] {
		t.Errorf("Expected pump heat to be %v, got %v", Config["primary_loop"]["pump_off_heat"], pl.PumpHeat())
	}
}

func TestPrimaryLoop_AdjustBoronConcentration(t *testing.T) {
	pl := NewPrimaryLoop("TestLoop-Pump", "The is a test.")
	sim := simworks.NewSimulator("Test Sim", "Test Sim")
	sim.AddComponent(pl)

	// turn on pump and establish boron concentration target
	sim.QueueEvent(NewEvent_PumpSwitch(true))
	boronTarget := 300.0
	eventToWatch := NewEvent_BoronConcentration(boronTarget)
	sim.QueueEvent(eventToWatch)

	if !eventToWatch.IsPending() {
		t.Errorf("Expected event to be pending, got %v", eventToWatch.Status)
	}

	sim.Step()
	sim.Step()

	if !eventToWatch.IsInProgress() {
		t.Errorf("Expected event to be in progress, got %v", eventToWatch.Status)
	}

	// let 2 hours pass
	sim.RunForABit(0, 2, 0, 0)

	// check that the boron concentration has changed
	if pl.BoronConcentration() != boronTarget {
		t.Errorf("Expected boron concentration to be %v, got %v", boronTarget, pl.BoronConcentration())
	}
	if !eventToWatch.IsComplete() {
		t.Errorf("Expected event to be complete, got %v", eventToWatch.Status)
	}
}

func TestPrimaryLoop_BoronConcentrationWhenPumpInitiallyOff(t *testing.T) {
	pl := NewPrimaryLoop("TestLoop-PumpOff", "This is a test.")
	sim := simworks.NewSimulator("Test Sim", "Test Sim")
	sim.AddComponent(pl)

	// Ensure pump is off
	if pl.PumpOn() {
		t.Errorf("Expected pump to be off initially, got %v", pl.PumpOn())
	}

	initialBoron := pl.BoronConcentration()

	// Try to adjust boron concentration
	boronTarget := 50.0
	boronEvent := NewEvent_BoronConcentration(boronTarget)
	sim.QueueEvent(boronEvent)
	sim.RunForABit(0, 0, 1, 0) // Run for 1 minute

	// Check that boron concentration hasn't changed
	if pl.BoronConcentration() != initialBoron {
		t.Errorf("Expected boron concentration to remain at %v, but got %v", initialBoron, pl.BoronConcentration())
	}

	// Verify that the event is still in progress (not completed)
	if !boronEvent.IsInProgress() {
		t.Errorf("Expected boron adjustment event to be in progress, but got status: %v", boronEvent.Status)
	}

	sim.QueueEvent(NewEvent_PumpSwitch(true))
	sim.Step()
	if !pl.PumpOn() {
		t.Errorf("Expected pump to be on, got %v", pl.PumpOn())
	}

	// now boron should start to change
	sim.Step()
	if pl.BoronConcentration() == initialBoron {
		t.Errorf("Expected boron concentration to change, but it remained at %v", initialBoron)
	}

	sim.RunForABit(0, 1, 0, 0)
	if pl.BoronConcentration() != boronTarget {
		t.Errorf("Expected boron to reach target, but it is at %v", pl.BoronConcentration())
	}

	if !boronEvent.IsComplete() {
		t.Errorf("Expected boron adjustment event to be complete, but got status: %v", boronEvent.Status)
	}
}
