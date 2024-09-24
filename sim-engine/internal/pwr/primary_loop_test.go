package pwr

import (
	"testing"

	"worldofnuclear.com/internal/simworks"
)

func TestNewPrimaryLoop(t *testing.T) {
	pl := NewPrimaryLoop("TestLoop-Pump", "The is a test.")
	if pl.Name() != "TestLoop-Pump" {
		t.Errorf("Expected name to be TestLoop-Pump, got %s", pl.Name())
	}
	if pl.Description() != "The is a test." {
		t.Errorf("Expected description to be The is a test., got %s", pl.Description())
	}
}

func TestSimulatorDrivesPrimaryLoop(t *testing.T) {
	pl := NewPrimaryLoop("TestLoop-Pump", "The is a test.")
	sim := simworks.NewSimulator("Test Sim", "Test Sim")
	sim.AddComponent(pl)
	sim.Run(100)

	if pl.LatestMoment().IsZero() {
		t.Errorf("Expected LatestMoment to be non-zero, got %v", pl.LatestMoment())
	}
}

// func TestPrimaryLoopPump(t *testing.T) {
// 	simulator := NewSimulator("Test Sim", "Early detection is the key.", "PWR")

// 	pl := sim.NewPrimaryLoop("TestLoop-Pump")
// 	testy, testEnv := setupSimulationEnvironment()
// 	testy.AddComponent(pl)

// 	// Initially, the pump should be off and pressure should be 0
// 	if pl.Pressure() != 0 {
// 		t.Errorf("Initial pump pressure should be 0, got %f", pl.Pressure())
// 	}
// 	if pl.FlowVolume() != 0 {
// 		t.Errorf("Initial flow volume should be 0, got %f", pl.FlowVolume())
// 	}

// 	// Turn on the pump
// 	pl.SwitchOnPump()

// 	// Update the primary loop to apply changes
// 	pl.Update(testEnv, testy)

// 	// Check that the pressure is greater than 0
// 	if pl.Pressure() <= 0 {
// 		t.Errorf("Pump pressure should be greater than 0 when pump is on, got %f", pl.Pressure())
// 	}

// 	if pl.FlowVolume() <= 0 {
// 		t.Errorf("Flow volume should be greater than 0 when pump is on, got %f", pl.FlowVolume())
// 	}

// 	pl.SwitchOffPump()
// 	pl.Update(testEnv, testy)

// 	if pl.Pressure() != 0 {
// 		t.Errorf("Pressure should return to 0, got %f", pl.Pressure())
// 	}
// 	if pl.FlowVolume() != 0 {
// 		t.Errorf("Flow volume should return to 0, got %f", pl.FlowVolume())
// 	}

// }

// func TestPrimaryLoopBoronAdjustments(t *testing.T) {
// 	pl := sim.NewPrimaryLoop("TestLoop-BoronAdjustments")
// 	testy, testEnv := setupSimulationEnvironment()
// 	testy.AddComponent(pl)

// 	// Initially, the boron concentration should be 0
// 	if pl.BoronConcentration() != 0 {
// 		t.Errorf("Initial boron concentration should be 0, got %f", pl.BoronConcentration())
// 	}

// 	// Set a target boron concentration
// 	pl.AdjustBoronConcentrationTarget(100)
// 	pl.SwitchOnPump()
// 	pl.Update(testEnv, testy)

// 	// Check that the boron concentration is greater than 0
// 	if pl.BoronConcentration() <= 0 {
// 		t.Errorf("Boron concentration should be greater than 0 after one tick when target is set, got %f", pl.BoronConcentration())
// 	}

// 	// run for just over an hour
// 	for i := 0; i < 62; i++ {
// 		pl.Update(testEnv, testy)
// 	}
// 	if pl.BoronConcentration() != 100 {
// 		t.Errorf("Boron concentration should have reached target by now, got %f", pl.BoronConcentration())
// 	}

// 	// raise it again
// 	pl.AdjustBoronConcentrationTarget(200)
// 	// run for about an hour
// 	for i := 0; i < 62; i++ {
// 		pl.Update(testEnv, testy)
// 	}
// 	if pl.BoronConcentration() != 200 {
// 		t.Errorf("Boron concentration should have reached target by now, got %f", pl.BoronConcentration())
// 	}

// 	// now lower it
// 	pl.AdjustBoronConcentrationTarget(150)
// 	// run for about half an hour
// 	for i := 0; i < 32; i++ {
// 		pl.Update(testEnv, testy)
// 	}
// 	if pl.BoronConcentration() != 150 {
// 		t.Errorf("Boron concentration should have dropeed to target by now, got %f", pl.BoronConcentration())
// 	}

// 	pl.SwitchOffPump()
// 	pl.AdjustBoronConcentrationTarget(50)
// 	pl.Update(testEnv, testy)

// 	if pl.BoronConcentration() != 150 {
// 		t.Errorf("Boron concentration should not change when pump is off, got %f", pl.BoronConcentration())
// 	}
// }
