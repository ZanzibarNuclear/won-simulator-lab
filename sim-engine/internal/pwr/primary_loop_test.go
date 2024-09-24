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

func TestPumpSwitch(t *testing.T) {
	pl := NewPrimaryLoop("TestLoop-Pump", "The is a test.")
	sim := simworks.NewSimulator("Test Sim", "Test Sim")
	sim.AddComponent(pl)
	sim.QueueEvent(simworks.NewImmediateEventBool(Event_pl_pumpSwitch, true))
	sim.Run(1)
	// pump should be running with side effects
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

	sim.QueueEvent(simworks.NewImmediateEventBool(Event_pl_pumpSwitch, false))
	sim.Run(1)
	// pump should be off, related values at low state
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

func TestAdjustBoronConcentration(t *testing.T) {
	pl := NewPrimaryLoop("TestLoop-Pump", "The is a test.")
	sim := simworks.NewSimulator("Test Sim", "Test Sim")
	sim.AddComponent(pl)

	// turn on pump and establish boron concentration target
	sim.QueueEvent(simworks.NewImmediateEventBool(Event_pl_pumpSwitch, true))
	boronEvent := simworks.NewAdjustmentEvent(Event_pl_boronConcentration, 300.0)
	sim.QueueEvent(boronEvent)
	eventToWatch := &sim.Events[1]

	if !eventToWatch.IsPending() {
		t.Errorf("Expected event to be pending, got %v", eventToWatch.Status)
	}

	sim.Run(2)

	if eventToWatch.Code != Event_pl_boronConcentration {
		t.Errorf("Expected event to be %v, got %v", Event_pl_boronConcentration, eventToWatch.Code)
	}
	if !eventToWatch.IsInProgress() {
		t.Errorf("Expected event to be in progress, got %v", eventToWatch.Status)
	}

	// let 2 hours pass
	sim.RunForABit(0, 2, 0, 0)

	// check that the boron concentration has changed
	if pl.BoronConcentration() != 300.0 {
		t.Errorf("Expected boron concentration to be %v, got %v", 300.0, pl.BoronConcentration())
	}
	if !eventToWatch.IsComplete() {
		t.Errorf("Expected event to be complete, got %v", eventToWatch.Status)
	}
}
