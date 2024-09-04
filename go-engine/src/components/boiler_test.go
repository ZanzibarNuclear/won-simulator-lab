package components_test

import (
	"testing"
	"won/sim-lab/go-engine/components"
)

func TestBoiler(t *testing.T) {
	boiler := components.NewBoiler()

	if boiler.Running() {
		t.Error("New boiler instance should be running by default")
	}

	boiler.TurnOn()

	if !boiler.Running() {
		t.Error("Boiler should be running after turning on")
	}
}