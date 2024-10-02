package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"worldofnuclear.com/internal/simworks"
)

func TestNewSteamGenerator(t *testing.T) {
	sg := NewSteamGenerator("SG1", "Test Steam Generator", nil, nil)
	s := simworks.NewSimulator("Test Simulator", "Testing steam generator")
	s.AddComponent(sg)

	roomTemp := Config["common"]["room_temperature"]

	assert.Equal(t, "SG1", sg.Name())
	assert.Equal(t, "Test Steam Generator", sg.Description())
	assert.Equal(t, roomTemp, sg.PrimaryInletTemp())
	assert.Equal(t, roomTemp, sg.PrimaryOutletTemp())
	assert.Equal(t, roomTemp, sg.SecondaryInletTemp())
	assert.Equal(t, roomTemp, sg.SecondaryOutletTemp())
	assert.Equal(t, 0.0, sg.HeatTransferRate())
	assert.Equal(t, 0.0, sg.SteamFlowRate())

	_, err := sg.Update(s)
	assert.Error(t, err)
}

func TestSteamGenerator_NormalConditions(t *testing.T) {
	primaryLoop := &PrimaryLoop{}
	secondaryLoop := &SecondaryLoop{}
	sg := NewSteamGenerator("SG1", "Test Steam Generator", primaryLoop, secondaryLoop)
	s := simworks.NewSimulator("Test Simulator", "Testing steam generator")
	s.AddComponent(primaryLoop)
	s.AddComponent(secondaryLoop)
	s.AddComponent(sg)

	// set conditions for primary loop
	primaryLoop.SetHotLegTemperature(325.0)

	// set conditions for secondary loop
	// TODO:

	s.Step()

	// TODO: check that primary inlet matches hot leg temperature
	// TODO: check cold leg temperature is between room temp and 325.0

	assert.Fail(t, "Not implemented")
}
