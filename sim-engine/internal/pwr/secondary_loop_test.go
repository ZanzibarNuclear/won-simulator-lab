package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"worldofnuclear.com/internal/simworks"
)

func TestSecondaryLoop_Init(t *testing.T) {
	sl := NewSecondaryLoop("Test Secondary Loop", "A test secondary loop")

	assert.NotNil(t, sl)
	assert.Equal(t, "Test Secondary Loop", sl.Name())
	assert.Equal(t, "A test secondary loop", sl.Description())

	// Check default values
	assert.Equal(t, Config["common"]["room_temperature"], sl.SteamTemperature())
	assert.Equal(t, Config["common"]["atmospheric_pressure"], sl.SteamPressure())
	assert.Equal(t, Config["secondary_loop"]["base_feedwater_temperature"], sl.FeedwaterTemperatureOut())
	assert.Equal(t, Config["secondary_loop"]["base_feedwater_temperature"], sl.FeedwaterTemperatureIn())

	// Check initial states
	assert.False(t, sl.PowerOperatedReliefValveOpen())
	assert.False(t, sl.FeedwaterPumpOn())
	assert.Equal(t, 0.0, sl.FeedwaterFlowRate())
	assert.False(t, sl.FeedheatersOn())
}

func TestSecondaryLoop_FeedwaterPump(t *testing.T) {
	sl := NewSecondaryLoop("Test Secondary Loop", "A test secondary loop")
	s := simworks.NewSimulator("Test Sim", "Test Sim")
	s.AddComponent(sl)
	s.Step()

	// Initially, the pump should be off
	assert.False(t, sl.FeedwaterPumpOn())
	assert.Equal(t, 0.0, sl.FeedwaterFlowRate())

	// Queue an event to turn the pump on
	s.QueueEvent(NewEvent_FeedwaterPumpSwitch(true))
	s.Step()

	assert.True(t, sl.FeedwaterPumpOn())
	assert.Greater(t, sl.FeedwaterFlowRate(), 0.0)

	s.RunForABit(0, 0, 1, 0)
	assert.Equal(t, Config["secondary_loop"]["feedwater_flow_rate_target"], sl.FeedwaterFlowRate())

	// Queue an event to turn the pump on
	s.QueueEvent(NewEvent_FeedwaterPumpSwitch(false))
	s.Step()

	// Turn the pump off
	assert.False(t, sl.FeedwaterPumpOn())
	assert.Less(t, sl.FeedwaterFlowRate(), Config["secondary_loop"]["feedwater_flow_rate_target"])

	s.RunForABit(0, 0, 1, 0)
	assert.Equal(t, 0.0, sl.FeedwaterFlowRate())
}

func TestSecondaryLoop_Feedheaters(t *testing.T) {
	sl := NewSecondaryLoop("Test Secondary Loop", "A test secondary loop")
	s := simworks.NewSimulator("Test Sim", "Test Sim")
	s.AddComponent(sl)
	s.Step()

	// Initially, the feedheaters should be off
	assert.False(t, sl.FeedheatersOn())
	assert.Equal(t, Config["secondary_loop"]["base_feedwater_temperature"], sl.FeedwaterTemperatureOut())

	s.QueueEvent(NewEvent_FeedheatersSwitch(true))
	s.Step()

	// Turn the feedheaters on
	sl.SwitchFeedheaters(true)
	assert.True(t, sl.FeedheatersOn())
	assert.Greater(t, sl.FeedwaterTemperatureOut(), Config["secondary_loop"]["base_feedwater_temperature"])

	s.RunForABit(0, 0, 3, 0)
	assert.Equal(t, Config["secondary_loop"]["heated_feedwater_temperature"], sl.FeedwaterTemperatureOut())

	s.QueueEvent(NewEvent_FeedheatersSwitch(false))
	s.Step()

	// Turn the feedheaters off
	sl.SwitchFeedheaters(false)
	assert.False(t, sl.FeedheatersOn())
	assert.Less(t, sl.FeedwaterTemperatureOut(), Config["secondary_loop"]["heated_feedwater_temperature"])

	s.RunForABit(0, 0, 3, 0)
	assert.Equal(t, Config["secondary_loop"]["base_feedwater_temperature"], sl.FeedwaterTemperatureOut())
}

func TestSecondaryLoop_PowerOperatedReliefValve(t *testing.T) {
	sl := NewSecondaryLoop("Test Secondary Loop", "Power Operated Relief Valve test")
	s := simworks.NewSimulator("Test Sim", "Test Sim")
	s.AddComponent(sl)

	// build steam pressure (or set it as a shortcut)

	sl.SetSteamPressure(5.0)
	s.Step()

	// open the value
	pressureMark := sl.SteamPressure()
	temperatureMark := sl.SteamTemperature()

	s.QueueEvent(NewEvent_PowerOperatedReliefValveSwitch(true))
	s.Step()
	assert.True(t, sl.PowerOperatedReliefValveOpen())

	s.RunForABit(0, 0, 0, 4)

	assert.Less(t, sl.SteamPressure(), pressureMark)
	assert.Less(t, sl.SteamTemperature(), temperatureMark)

	// close the value - pressure and temperature should hold
	pressureMark = sl.SteamPressure()
	temperatureMark = sl.SteamTemperature()

	s.QueueEvent(NewEvent_PowerOperatedReliefValveSwitch(false))
	s.Step()

	assert.False(t, sl.PowerOperatedReliefValveOpen())
	assert.Equal(t, pressureMark, sl.SteamPressure())
	assert.Equal(t, temperatureMark, sl.SteamTemperature())
}

func TestSecondaryLoop_MssvVentEvent(t *testing.T) {
	sl := NewSecondaryLoop("Test Secondary Loop", "A test secondary loop")
	s := NewPwrSim("Test PWR", "Relief valve vent test simulation")
	s.AddComponent(sl)
	s.SetEventHandler(s)

	s.Step()

	excessivePressure := Config["secondary_loop"]["mssv_pressure_threshold"] + 0.1
	sl.SetSteamPressure(excessivePressure)
	s.Step()

	// Find the relief valve event in InactiveEvents
	var mssvVentEvent *simworks.Event
	for _, event := range s.InactiveEvents {
		// t.Logf("Inactive event: %v", event)
		if event.Code == Event_sl_emergencyMssvVent {
			mssvVentEvent = event
			break
		}
	}

	if mssvVentEvent == nil {
		t.Error("Relief valve vent event not found in InactiveEvents")
		return
	}
	if !mssvVentEvent.IsComplete() {
		t.Error("Relief valve vent event should have been processed by simulator")
	}

	assert.Less(t, sl.SteamPressure(), excessivePressure)
}
