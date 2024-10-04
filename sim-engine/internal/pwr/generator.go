package pwr

import (
	"errors"
	"fmt"

	"worldofnuclear.com/internal/simworks"
)

type Generator struct {
	simworks.BaseComponent
	rpm                int
	electricalPowerOut float64 // in megawatts (MW)
	connectedToGrid    bool
	poles              int     // 2 for 2-pole, 4 for 4-pole
	standardFrequency  float64 // Hz; 60 in North America, 50 in Europe
	steamTurbine       *SteamTurbine
}

func NewGenerator(name string, description string, turbine *SteamTurbine) *Generator {
	return &Generator{
		BaseComponent:      *simworks.NewBaseComponent(name, description),
		rpm:                0,
		electricalPowerOut: 0,
		connectedToGrid:    false,
		steamTurbine:       turbine,
		poles:              2,
		standardFrequency:  60.0,
	}
}

func (g *Generator) Rpm() int {
	return g.rpm
}

func (g *Generator) ElectricalPowerOut() float64 {
	return g.electricalPowerOut
}

func (g *Generator) ConnectedToGrid() bool {
	return g.connectedToGrid
}

func (g *Generator) Status() map[string]interface{} {
	return map[string]interface{}{
		"about":              g.BaseComponent.Status(),
		"rpm":                g.rpm,
		"connectedToGrid":    g.connectedToGrid,
		"electricalPowerOut": g.electricalPowerOut,
	}
}

func (g *Generator) Print() {
	fmt.Printf("=> Electrical Generator\n")
	g.BaseComponent.Print()
	fmt.Printf("RPM: %d\n", g.rpm)
	fmt.Printf("Connected To Grid: %t\n", g.connectedToGrid)
	fmt.Printf("Electrical Power Out: %.2f MW\n", g.electricalPowerOut)
}

func (g *Generator) Update(s *simworks.Simulator) (map[string]interface{}, error) {
	if g.steamTurbine == nil {
		return g.Status(), errors.New("steam turbine not connected")
	}

	g.BaseComponent.Update(s)
	for _, event := range s.Events {
		if event.IsInProgress() {
			g.processEvent(event)
		}
	}

	g.rpm = g.steamTurbine.Rpm()
	if g.connectedToGrid {
		g.electricalPowerOut = float64(g.rpm) * 0.001 // Arbitrary scaling factor??
	} else {
		g.electricalPowerOut = 0
	}

	// once attached to grid, fall out if frequency (RPMs) changes too much
	// frequency must stay within 0.067 Hz (about 0.1% of 60 Hz)
	frequencyTolerance := g.standardFrequency * 0.001
	if g.connectedToGrid {
		if !simworks.AlmostEqual(float64(g.rpm), g.standardFrequency, frequencyTolerance) {
			g.connectedToGrid = false
			return g.Status(), errors.New("frequency changed too much")
		}
		// also depends on phase angle and voltage; not modeled as yet
	}

	return g.Status(), nil
}

func (g *Generator) processEvent(event *simworks.Event) {
	switch event.Code {
	case Event_g_connectToGrid:

		// frequency has to match standard frequency
		// TODO: don't just give up; tell steam turbine to go faster or slower to match frequency
		frequency := 1.0 / float64(g.rpm)
		if !simworks.AlmostEqual(frequency, Config["generator"]["standard_ac_frequency"], 0.1) {
			event.SetCanceled() // TODO: add cancelation reason to event
			return
		}
		g.connectedToGrid = event.Truthy()
		event.SetComplete()
	}
}
