package sim

import (
	"testing"
)

// setupSimulationEnvironment creates and returns a new Simulation and Environment
func setupSimulationEnvironment() (*Simulation, *Environment) {
	simulation := NewSimulation("Test Sim", "Early detection is the key.")
	environment := NewEnvironment()
	return simulation, environment
}

func TestSecondaryLoopInitialConditions(t *testing.T) {
	// Create a new secondary loop
	sl := NewSecondaryLoop("test")

	// Check initial conditions
	if sl.feedwaterPumpOn {
		t.Errorf("Feedwater pump should be off initially, got %t", sl.feedwaterPumpOn)
	}

	if sl.feedheatersOn {
		t.Errorf("Feedheaters should be off initially, got %t", sl.feedheatersOn)
	}

	sl.SwitchOnFeedwaterPump()
	sl.SwitchOnFeedheaters()

	mySim, myEnv := setupSimulationEnvironment()
	mySim.AddComponent(sl)

	sl.Update(myEnv, mySim)

}

func TestSecondaryLoopFeedwaterPump(t *testing.T) {
	// Create a new secondary loop
	sl := NewSecondaryLoop("TestLoop-Feedwater")

	// Create a new simulation and environment
	sim, env := setupSimulationEnvironment()
	sim.AddComponent(sl)

	// Run update once
	sl.Update(env, sim)

	// Check that FeedwaterVolume is 0 initially
	if sl.FeedwaterVolume() != 0 {
		t.Errorf("Initial feedwater volume should be 0, got %f", sl.FeedwaterVolume())
	}

	// Turn the pump on
	sl.SwitchOnFeedwaterPump()
	sl.Update(env, sim)

	// Check that FeedwaterVolume is > 0 after pump is turned on
	if sl.FeedwaterVolume() <= 0 {
		t.Errorf("Feedwater volume should be greater than 0 when pump is on, got %f", sl.FeedwaterVolume())
	}

	sl.SwitchOffFeedwaterPump()
	sl.Update(env, sim)

	// Check that FeedwaterVolume returns to 0
	if sl.FeedwaterVolume() != 0 {
		t.Errorf("Feedwater volume should return to 0 when pump is off, got %f", sl.FeedwaterVolume())
	}
}

func TestSecondaryLoopFeedheaters(t *testing.T) {
	sl := NewSecondaryLoop("TestLoop-Feedheaters")
	sim, env := setupSimulationEnvironment()
	sim.AddComponent(sl)

	// Turn on the pump and leave it on for the whole test
	sl.SwitchOnFeedwaterPump()
	sl.Update(env, sim)

	// Check initial feedwater temperature
	initialTemp := sl.feedwaterTemperature
	if initialTemp != BASE_FEEDWATER_TEMPERATURE {
		t.Errorf("Initial feedwater temperature should be %f, got %f", BASE_FEEDWATER_TEMPERATURE, initialTemp)
	}

	// Turn on feedheaters
	sl.SwitchOnFeedheaters()

	// Run update multiple times to allow temperature to increase
	sl.Update(env, sim)

	// Check that feedwater temperature has increased
	if sl.feedwaterTemperature == initialTemp {
		t.Errorf("Feedwater temperature should have increased. Got %f, expected > %f", sl.feedwaterTemperature, initialTemp)
	}

	// Run update multiple times to allow temperature to increase
	for i := 0; i < 30; i++ {
		sl.Update(env, sim)
	}

	// Check that feedwater temperature is at or approaching target temperature
	targetTemp := sl.TargetFeedwaterTemperature()
	if sl.feedwaterTemperature < targetTemp {
		t.Errorf("Feedwater temperature have reached target temperature. Got %f, expected close to %f", sl.feedwaterTemperature, targetTemp)
	}

	// Turn off feedheaters
	sl.SwitchOffFeedheaters()
	sl.Update(env, sim)

	// Check that feedwater temperature has decreased
	cooledTemp := sl.feedwaterTemperature
	if cooledTemp >= targetTemp {
		t.Errorf("Feedwater temperature should have dropped a little. Got %f, expected < %f", cooledTemp, targetTemp)
	}

	// Run update multiple times to allow temperature to decrease
	for i := 0; i < 30; i++ {
		sl.Update(env, sim)
	}

	// Check that feedwater temperature is approaching base temperature
	if sl.feedwaterTemperature > BASE_FEEDWATER_TEMPERATURE {
		t.Errorf("Feedwater temperature should be at base temperature. Got %f, expected close to %f", sl.feedwaterTemperature, BASE_FEEDWATER_TEMPERATURE)
	}
}
