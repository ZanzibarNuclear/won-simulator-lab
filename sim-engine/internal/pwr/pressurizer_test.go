package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"worldofnuclear.com/internal/simworks"
)

func TestNewPressurizerDefaults(t *testing.T) {
	pl := NewPrimaryLoop("TestPrimaryLoop", "Test primary loop description")
	p := NewPressurizer("TestPressurizer", "Test pressurizer description", pl)

	// Test default values
	assert.NotEmpty(t, p.ID())
	assert.Equal(t, "TestPressurizer", p.Name())
	assert.Equal(t, "Test pressurizer description", p.Description())
	assert.Equal(t, Config["common"]["atmospheric_pressure"], p.Pressure())
	assert.Equal(t, Config["common"]["room_temperature"], p.Temperature())
	assert.False(t, p.HeaterOn())
	assert.False(t, p.SprayNozzleOpen())
}

func TestPressurizer_HeaterPowerEvent(t *testing.T) {
	pl := NewPrimaryLoop("TestPrimaryLoop", "Test primary loop description")
	p := NewPressurizer("TestPressurizer", "Test pressurizer description", pl)
	pwrSim := NewPwrSim("Test PWR", "Heater power test simulation")
	pwrSim.AddComponent(p)

	// Initially, the heater should be off
	assert.False(t, p.HeaterOn())

	// Create and queue the event to turn the heater on
	eventToWatch := NewEvent_HeaterPower(true)
	pwrSim.QueueEvent(eventToWatch)
	pwrSim.Step()

	assert.True(t, p.HeaterOn())
	assert.True(t, eventToWatch.IsComplete())

	// Now test turning the heater off
	eventToWatch = NewEvent_HeaterPower(false)
	pwrSim.QueueEvent(eventToWatch)
	pwrSim.Step()

	assert.False(t, p.HeaterOn())
	assert.True(t, eventToWatch.IsComplete())
}

func TestPressurizer_SprayNozzleEvent(t *testing.T) {
	pl := NewPrimaryLoop("TestPrimaryLoop", "Test primary loop description")
	p := NewPressurizer("TestPressurizer", "Test pressurizer description", pl)
	pwrSim := NewPwrSim("Test PWR", "Spray nozzle test simulation")
	pwrSim.AddComponent(p)

	// Initially, the spray nozzle should be closed
	assert.False(t, p.SprayNozzleOpen())
	assert.Equal(t, 0.0, p.SprayFlowRate())

	// Create and queue the event to open the spray nozzle
	eventToWatch := NewEvent_SprayNozzle(true)
	pwrSim.QueueEvent(eventToWatch)
	pwrSim.Step()

	assert.True(t, p.SprayNozzleOpen())
	assert.Equal(t, Config["pressurizer"]["spray_flow_rate"], p.SprayFlowRate())
	assert.True(t, eventToWatch.IsComplete())

	// Now test closing the spray nozzle
	if !eventToWatch.IsComplete() {
		t.Error("SprayNozzleEvent should be completed after one step")
	}

	// Check if spray flow rate is set correctly
	expectedSprayFlowRate := Config["pressurizer"]["spray_flow_rate"]
	assert.Equal(t, expectedSprayFlowRate, p.SprayFlowRate())

	// Now test closing the spray nozzle
	eventToWatch = NewEvent_SprayNozzle(false)
	pwrSim.QueueEvent(eventToWatch)
	pwrSim.Step()

	assert.False(t, p.SprayNozzleOpen())
	assert.Equal(t, 0.0, p.SprayFlowRate())
	assert.True(t, eventToWatch.IsComplete())
}

func TestPressurizer_SetPressureEvent(t *testing.T) {
	pl := NewPrimaryLoop("TestPrimaryLoop", "Test primary loop description")
	p := NewPressurizer("TestPressurizer", "Test pressurizer description", pl)
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

	assert.NotEqual(t, initialPressure, p.Pressure())
	assert.True(t, p.Pressure() >= targetPressure)
	assert.True(t, eventToWatch.IsComplete())
	assert.True(t, p.HeaterOnLow())
}

func TestPressurizer_ReliefValveVentEvent(t *testing.T) {
	pl := NewPrimaryLoop("TestPrimaryLoop", "Test primary loop description")
	p := NewPressurizer("TestPressurizer", "Test pressurizer description", pl)
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

	assert.True(t, pressureEvent.IsInProgress())
	assert.True(t, heaterEvent.IsComplete())

	pwrSim.RunForABit(0, 0, 0, 5)

	assert.Equal(t, 1, len(pwrSim.Events))
	assert.Equal(t, 1, len(pwrSim.InactiveEvents))

	pwrSim.RunForABit(0, 0, 1, 2)

	assert.Equal(t, 1, len(pwrSim.Events))
	assert.Equal(t, 2, len(pwrSim.InactiveEvents))

	// Find the relief valve event in InactiveEvents
	var reliefValveEvent *simworks.Event
	for _, event := range pwrSim.InactiveEvents {
		// t.Logf("Inactive event: %v", event)
		if event.Code == Event_pr_reliefValveVent {
			reliefValveEvent = event
			break
		}
	}

	assert.NotNil(t, reliefValveEvent)
	assert.True(t, reliefValveEvent.IsComplete())
}
