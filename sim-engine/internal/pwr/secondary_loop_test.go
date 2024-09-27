package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSecondaryLoop(t *testing.T) {
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
