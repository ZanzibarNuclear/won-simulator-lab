package sim

import (
	"fmt"
	"math"
)

type ReactorCore struct {
	BaseComponent
	reactivity            float64 // negative means subcritical, 0 means critical, positive means supercritical
	neutronFlux           float64 // in neutrons per second
	temperature           float64 // in degrees Celsius
	heatEnergyRate        float64 // in MW
	controlRods           *ControlRods
	primaryLoop           *PrimaryLoop
	withdrawShutdownBanks bool
	bringToCritical       bool // withdraw control rods until supercritical, then insert just a bit to get to critical
	scram                 bool
}

func NewReactorCore(name string) *ReactorCore {
	return &ReactorCore{
		BaseComponent:  BaseComponent{Name: name},
		reactivity:     0.0,
		neutronFlux:    0.0,
		temperature:    20.0, // Start at room temperature (Celsius)
		heatEnergyRate: 0.0,
		controlRods:    NewControlRods(),
	}
}

func (rc *ReactorCore) ConnectToPrimaryLoop(primaryLoop *PrimaryLoop) {
	rc.primaryLoop = primaryLoop
}

func (rc *ReactorCore) Update(env *Environment, s *Simulation) {
	// Update reactivity based on control rod insertion
	rc.reactivity = 1.0 - rc.controlRods.CalculateAverageInsertion()

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

func (rc *ReactorCore) BringToCritical() {
	rc.bringToCritical = true
}

func (rc *ReactorCore) InsertShutdownBanks() {
	rc.withdrawShutdownBanks = false
}

func (rc *ReactorCore) Scram() {
	rc.scram = true
}
