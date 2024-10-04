package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"worldofnuclear.com/internal/simworks"
)

func TestNewSteamTurbine(t *testing.T) {
	turbine := NewSteamTurbine("Test Turbine", "A test steam turbine", nil)
	s := simworks.NewSimulator("Test", "Test")
	s.AddComponent(turbine)

	assert.Equal(t, "Test Turbine", turbine.Name())
	assert.Equal(t, "A test steam turbine", turbine.Description())

	_, err := turbine.Update(s)
	assert.Error(t, err)
}

func TestSteamTurbine_NormalOperation(t *testing.T) {
	sl := NewSecondaryLoop("Test Secondary Loop", "A test secondary loop")
	turbine := NewSteamTurbine("Test Turbine", "A test steam turbine", sl)
	s := simworks.NewSimulator("Test", "Test")
	s.AddComponent(sl)
	s.AddComponent(turbine)

	// set steam pressure in secondary loop; find out RPMs of steam turbine.
	sl.SetSteamPressure(8.0)
	// raising steam pressure should raise RPMs
	s.RunForABit(0, 0, 1, 0)

	assert.Equal(t, 3600, turbine.Rpm())

	sl.Print()
	turbine.Print()

	// tune so that expected pressure produces sufficient RPMs to generate electricity that is suitable for grid

}
