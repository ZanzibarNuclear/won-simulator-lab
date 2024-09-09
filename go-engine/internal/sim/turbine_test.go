package sim_test

import (
	"testing"
	"won/sim-lab/go-engine/internal/sim"
)

func TestTurbine(t *testing.T) {
	turbine := sim.NewTurbine("test-turbine-1")

	if turbine.Rpm() != 0 {
		t.Error("New turbine instance should not be spinning")
	}

	if turbine.MaxedOut() {
		t.Error("New turbine instance should not be maxed out")
	}
}
