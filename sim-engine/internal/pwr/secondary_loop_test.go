package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"worldofnuclear.com/internal/simworks"
)

func TestNewSecondaryLoop_InitialState(t *testing.T) {
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
	sl := NewSecondaryLoop("Test Secondary Loop", "A test secondary loop")
	s := NewPwrSim("Test PWR", "Power Operated Relief Valve test simulation")
	s.AddComponent(sl)
	s.SetEventHandler(s)

	// Initially, the Power Operated Relief Valve should be closed
	assert.False(t, sl.PowerOperatedReliefValveOpen())

	// Queue an event to open the Power Operated Relief Valve
	s.QueueEvent(NewEvent_PowerOperatedReliefValveSwitch(true))
	s.Step()

	// Check if the valve is open
	assert.True(t, sl.PowerOperatedReliefValveOpen())

	// Queue an event to close the Power Operated Relief Valve
	s.QueueEvent(NewEvent_PowerOperatedReliefValveSwitch(false))
	s.Step()

	// Check if the valve is closed
	assert.False(t, sl.PowerOperatedReliefValveOpen())

	// Set steam pressure above the threshold
	sl.steamPressure = Config["secondary_loop"]["porv_pressure_threshold"] + 1.0
	s.Step()

	// Check if the valve automatically opens
	assert.True(t, sl.PowerOperatedReliefValveOpen())

	// Check if the steam pressure decreases
	assert.Less(t, sl.SteamPressure(), Config["secondary_loop"]["porv_pressure_threshold"]+1.0)

	// Run for a bit to let the pressure stabilize
	s.RunForABit(0, 0, 1, 0)

	// Check if the valve automatically closes when pressure is below threshold
	assert.False(t, sl.PowerOperatedReliefValveOpen())
	assert.LessOrEqual(t, sl.SteamPressure(), Config["secondary_loop"]["porv_pressure_threshold"])
}

func TestSecondaryLoop_MssvVentEvent(t *testing.T) {
	sl := NewSecondaryLoop("Test Secondary Loop", "A test secondary loop")
	s := simworks.NewSimulator("Test Sim", "Test Sim")
	s.AddComponent(sl)
	s.Step()

	// Initially, the MSSV should be closed
	assert.False(t, sl.PowerOperatedReliefValveOpen())

	sl.steamPressure = Config["secondary_loop"]["mssv_pressure_threshold"] + 0.1

	// Check if the MSSV opens
	assert.True(t, sl.PowerOperatedReliefValveOpen())

	// Check if the steam pressure decreases
	assert.Less(t, sl.SteamPressure(), Config["secondary_loop"]["mssv_pressure_threshold"]+0.1)

	// Find the relief valve event in InactiveEvents
	var mssvVentEvent *simworks.Event
	for _, event := range s.InactiveEvents {
		// t.Logf("Inactive event: %v", event)
		if event.Code == Event_sl_emergencyMSSVReleased {
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

}
