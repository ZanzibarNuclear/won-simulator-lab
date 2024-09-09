package sim_test

import (
	"testing"
	"won/sim-lab/go-engine/internal/sim"
)

func TestBoiler(t *testing.T) {
	boiler := sim.NewBoiler("test-boiler-1")

	if boiler.Running() {
		t.Error("New boiler instance should be running by default")
	}

	boiler.TurnOn()

	if !boiler.Running() {
		t.Error("Boiler should be running after turning on")
	}
}