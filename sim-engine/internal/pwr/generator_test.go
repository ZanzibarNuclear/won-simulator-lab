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

func TestGenerator_IntegrationTest(t *testing.T) {
	// this looks like an integration test -- not appropriate for unit testing
	s := simworks.NewSimulator("Test Simulator", "A test simulator")
	pl := NewPrimaryLoop("Test Primary Loop", "A test primary loop")
	sl := NewSecondaryLoop("Test Secondary Loop", "A test secondary loop")
	sg := NewSteamGenerator("Test Steam Generator", "A test steam generator", pl, sl)
	st := NewSteamTurbine("Test Steam Turbine", "A test steam turbine", sl)
	g := NewGenerator("Test Generator", "A test generator", st)
	s.AddComponent(pl)
	s.AddComponent(sl)
	s.AddComponent(sg)
	s.AddComponent(st)
	s.AddComponent(g)
	s.Step()

	s.PrintStatus()

	assert.Fail(t, "Integration test; not implemented")
}

func TestGenerator_SteamTurbineDeterminesRPMs(t *testing.T) {
	// this looks like an integration test -- not appropriate for unit testing
	s := simworks.NewSimulator("Test Simulator", "A test simulator")
	st := &SteamTurbine{
		BaseComponent: *simworks.NewBaseComponent("Test Steam Turbine", "A test steam turbine"),
		secondaryLoop: &SecondaryLoop{
			BaseComponent: *simworks.NewBaseComponent("Test Secondary Loop", "A test secondary loop"),
			steamPressure: 8.0,
		},
	}
	g := NewGenerator("Test Generator", "A test generator", st)
	s.AddComponent(st)
	s.AddComponent(g)
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
