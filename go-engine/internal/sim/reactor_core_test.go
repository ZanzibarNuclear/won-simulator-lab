package sim

import (
	"fmt"
	"testing"
)

func setUpSimulation(boron int) *Simulation {
	sim := NewSimulation("Sim for Reactor Core Testing", "Customers do not replace QA.")

	// Create and add a primary loop with the right boron level for testing
	// primaryLoop := &PrimaryLoop{
	// 	BaseComponent: BaseComponent{
	// 		Name: "Test Primary Loop",
	// 	},
	// 	boronConcentration:       float64(boron),
	// 	boronConcentrationTarget: float64(boron),
	// 	pumpOn:                   true,
	// }
	primaryLoop := NewPrimaryLoop("Test Primary Loop")
	primaryLoop.AdjustBoronConcentrationTarget(float64(boron))
	primaryLoop.SwitchOnPump()
	sim.AddComponent(primaryLoop)

	// Create a reactor core and connect it to the primary loop
	reactorCore := NewReactorCore("Test Reactor Core")
	reactorCore.ConnectToPrimaryLoop(primaryLoop)
	sim.AddComponent(reactorCore)

	sim.Advance(50)

	if reactorCore.primaryLoop.boronConcentration != float64(boron) {
		fmt.Println("wth, with the primary loop?")
	}
	return sim
}

func TestInitialReactivityIsNegative(t *testing.T) {
	simulation := setUpSimulation(1900)

	reactorCore := simulation.FindReactorCore()
	primaryLoop := simulation.FindPrimaryLoop()

	if primaryLoop.boronConcentration != 1900 {
		t.Errorf("Expected boron concentration to be 1900, got: %f", primaryLoop.boronConcentration)
	}

	// Update the reactor core once to ensure initial values are set
	simulation.Advance(1)

	// Check if the initial reactivity is negative
	if reactorCore.reactivity >= 0 {
		t.Errorf("Initial reactivity should be negative, got: %f", reactorCore.reactivity)
	}
}
