package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"worldofnuclear.com/internal/simworks"
)

func TestCondenser_Init(t *testing.T) {
	c := NewCondenser("Test Condenser", "A test condenser", nil, nil)
	s := simworks.NewSimulator("Test Simulator", "A test simulator")
	s.AddComponent(c)

	assert.Equal(t, c.TubeMaterial(), "titanium")
	assert.Equal(t, c.CondenserType(), "water-cooled")
	assert.Equal(t, c.SurfaceArea(), Config["condenser"]["surface_area"])

	_, err := c.Update(s)
	assert.Error(t, err)
}

func TestCondenser_NormalOperation(t *testing.T) {
	// Create a mock steam turbine
	sl := NewSecondaryLoop("Test Secondary Loop", "A test secondary loop")
	st := &SteamTurbine{
		BaseComponent: *simworks.NewBaseComponent("Test Steam Turbine", "A test steam turbine"),
		thermalPower:  2000.0,
		secondaryLoop: sl,
	}
	c := NewCondenser("Test Condenser", "A test condenser", st, sl)

	s := simworks.NewSimulator("Test Simulator", "A test simulator")
	s.AddComponent(sl)
	s.AddComponent(st)
	s.AddComponent(c)

	// Run the update method to simulate normal operation
	s.Step()

	assert.Greater(t, c.HeatRejection(), 0.0) // starting quantity of heat, derived from steam turbine

	assert.Equal(t, c.CoolingWaterTempIn(), Config["condenser"]["cooling_water_temp_in"])
	assert.Greater(t, c.CoolingWaterTempOut(), c.CoolingWaterTempIn()) // temperature should rise because heat is being extracted
}
