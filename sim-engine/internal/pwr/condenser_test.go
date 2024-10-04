package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"worldofnuclear.com/internal/simworks"
)

func TestNewCondenser(t *testing.T) {
	c := NewCondenser("Test Condenser", "A test condenser", nil, nil)
	s := simworks.NewSimulator("Test Simulator", "A test simulator")
	s.AddComponent(c)

	_, err := c.Update(s)
	assert.Error(t, err)
}

func TestCondenserNormalOperation(t *testing.T) {
	// Create a mock steam turbine
	mockSteamTurbine := &SteamTurbine{
		thermalPower: 3000, // MW
	}
	sl := NewSecondaryLoop("Test Secondary Loop", "A test secondary loop")

	// Create a new condenser with the mock steam turbine
	c := NewCondenser("Test Condenser", "A test condenser", mockSteamTurbine, sl)
	s := simworks.NewSimulator("Test Simulator", "A test simulator")
	s.AddComponent(c)

	// Set initial conditions
	c.coolingWaterTempIn = 25.0 // °C

	// Run the update method to simulate normal operation
	_, err := c.Update(s)
	assert.NoError(t, err)

	assert.Fail(t, "Test not implemented")
}
