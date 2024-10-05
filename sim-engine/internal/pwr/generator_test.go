package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"worldofnuclear.com/internal/simworks"
)

func setupGenerator() (*Generator, *simworks.Simulator) {
	s := simworks.NewSimulator("Test Simulator", "A test simulator")
	sl := NewSecondaryLoop("Test Secondary Loop", "A test secondary loop")
	st := NewSteamTurbine("Test Steam Turbine", "A test steam turbine", sl)
	g := NewGenerator("Test Generator", "A test generator", st)
	s.AddComponent(sl)
	s.AddComponent(st)
	s.AddComponent(g)

	return g, s
}

func TestGenerator_InitWithoutSteamTurbine(t *testing.T) {
	g := NewGenerator("Test Generator", "A test generator", nil)
	s := simworks.NewSimulator("Test Simulator", "A test simulator")
	s.AddComponent(g)

	_, err := g.Update(s)
	assert.Error(t, err) // steam turbine is required
}

func TestGenerator_SteamTurbineDeterminesRPMs(t *testing.T) {
	g, s := setupGenerator()
	s.Step()

	assert.Equal(t, 3600, g.Rpm())
}

func TestGenerator_ElectricalPowerOut(t *testing.T) {
	g, s := setupGenerator()
	s.Step()

	assert.Equal(t, 2000.0, g.ElectricalPowerOut())
}

func TestGenerator_ConnectedToGrid(t *testing.T) {
	g, s := setupGenerator()
	s.Step()

	assert.False(t, g.ConnectedToGrid())

	assert.Fail(t, "Success case not implemented")
}
