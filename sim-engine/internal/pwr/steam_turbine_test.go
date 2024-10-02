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
	assert.Fail(t, "Not implemented")
}
