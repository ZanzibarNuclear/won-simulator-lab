package sim

import (
	"fmt"
	"math"
)

type ReactorCore struct {
	BaseComponent
	fuelAge               int     // in minutes
	reactivity            float64 // negative means subcritical, 0 means critical, positive means supercritical
	neutronFlux           float64 // in neutrons per second
	temperature           float64 // in degrees Celsius
	heatEnergyRate        float64 // in MW
	controlRods           *ControlRods
	primaryLoop           *PrimaryLoop
	withdrawShutdownBanks bool
	scram                 bool
}

func NewReactorCore(name string) *ReactorCore {
	return &ReactorCore{
		BaseComponent:  BaseComponent{Name: name},
		fuelAge:        0, // start w/ brand new fuel; this is something to play with, roll a die to pick a starting age, or let the user specify
		reactivity:     -1.0,
		neutronFlux:    0.1,
		temperature:    20.0, // Start at room temperature (Celsius)
		heatEnergyRate: 0.0,
		controlRods:    NewControlRods(),
	}
}

func (rc *ReactorCore) ConnectToControlRods(controlRods *ControlRods) {
	rc.controlRods = controlRods
}

func (rc *ReactorCore) ConnectToPrimaryLoop(primaryLoop *PrimaryLoop) {
	rc.primaryLoop = primaryLoop
}

const (
	BEGINNING_OF_CYCLE = WEEK_OF_MINUTES * 6
	MIDDLE_OF_CYCLE    = WEEK_OF_MINUTES * 66
	END_OF_CYCLE       = WEEK_OF_MINUTES * 6
)

func lookupPartOfCycle(fuelAge int) string {
	switch {
	case fuelAge < BEGINNING_OF_CYCLE:
		return "beginning"
	case fuelAge < MIDDLE_OF_CYCLE:
		return "middle"
	case fuelAge < END_OF_CYCLE:
		return "end"
	default:
		return "refueling"
	}
}

// use age of fuel to determine appropriate levels of boron concentration and control rod extraction
func lookupLevels(fuelAge int) (boronConcentration, controlRodWithdrawal float64) {
	switch partOfCycle := lookupPartOfCycle(fuelAge); partOfCycle {
	case "beginning":
		return 1900, 0.98
	case "middle":
		return 900, 0.8
	case "end":
		return 30, 0.5
	case "refueling":
		return 2200, 50
	}
	return 0, 0
}

/*
Young fuel (early in the cycle) requires more boron to moderate well.

	  Max: 2500
	  Beginning: 1200 - 2000 ppm
		Middle: 800 - 1000 ppm
		End: 10 - 50 ppm
		Refueling: 2000 - 2200 ppm

Notes: should be reduced gradually as fuel depletes
*/
func (rc *ReactorCore) Update(env *Environment, s *Simulation) {
	rc.fuelAge++ // keep fuel age in sync with sim time; TODO: improve by basing on operational minutes, not just elapsed time

	if rc.scram {
		rc.controlRods.Scram()
	}

	// determine reactivity

	// factors that affect reactivity
	// 1) PrimaryLoop.boronConcentration
	// 2) ControlRods.position: Shutdown rods out, control rods partially out
	// 3) Place in fuel cycle

	targetBoronConcentration, targetControlRodExtraction := lookupLevels(rc.fuelAge)

	// formula: ρ = (k - 1) / k, where ρ is reactivity and k is the effective multiplication factor

	// key assumption for the model:
	//   non-zero boron concentration is at a level where core goes critical with control rods half-way out
	//   this is a major simplification, but ought to be enough detail to tease out the interaction

	// sum up reactivity from all factors; TODO: find a more predictive formula
	reactivity := 0.0
	if rc.primaryLoop.boronConcentration < targetBoronConcentration {
		reactivity += 0.01
	} else if rc.primaryLoop.boronConcentration > targetBoronConcentration {
		reactivity -= 0.01
	}
	if rc.controlRods.AverageControlRodExtraction() < targetControlRodExtraction {
		reactivity -= 0.01
	} else if rc.controlRods.AverageControlRodExtraction() > targetControlRodExtraction {
		reactivity += 0.01
	}
	rc.reactivity = reactivity

	// use reactivity to determine flux
	if rc.reactivity > 0 {
		rc.neutronFlux *= 1.1
	} else if rc.reactivity < 0 {
		rc.neutronFlux *= 0.1
	}

	// Update heat energy rate based on reactivity
	rc.heatEnergyRate = 3000.0 * rc.reactivity // Assuming max output of 3000 MW

	// Simple temperature model (this should be more complex in reality)
	rc.temperature += (rc.heatEnergyRate / 1000.0) * 0.1          // Simplified heating
	rc.temperature = math.Max(20, math.Min(rc.temperature, 1000)) // Limit temperature range
}

func (rc *ReactorCore) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":           rc.Name,
		"reactivity":     rc.reactivity,
		"neutronFlux":    rc.neutronFlux,
		"temperature":    rc.temperature,
		"heatEnergyRate": rc.heatEnergyRate,
		"controlRods":    rc.controlRods,
	}
}

func (rc *ReactorCore) PrintStatus() {
	fmt.Printf("Reactor Core: %s\n", rc.Name)
	fmt.Printf("\tReactivity: %.2f\n", rc.reactivity)
	fmt.Printf("\tNeutron Flux: %.2f\n", rc.neutronFlux)
	fmt.Printf("\tTemperature: %.2f°C\n", rc.temperature)
	fmt.Printf("\tHeat Energy Rate: %.2f MW\n", rc.heatEnergyRate)
	fmt.Printf("\tControl Rods: %v\n", rc.controlRods)
}

func (rc *ReactorCore) HeatEnergyRate() float64 {
	return rc.heatEnergyRate
}

func (rc *ReactorCore) WithdrawShutdownBanks() {
	rc.withdrawShutdownBanks = true
}

func (rc *ReactorCore) InsertShutdownBanks() {
	rc.withdrawShutdownBanks = false
}

func (rc *ReactorCore) ScramReactor() {
	rc.scram = true
}

func (rc *ReactorCore) CancelScram() {
	rc.scram = false
}
