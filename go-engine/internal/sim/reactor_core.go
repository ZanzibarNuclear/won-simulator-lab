package sim

import (
	"fmt"
	"math"
)

type ReactorCore struct {
	BaseComponent
	controlRodInsertion float64 // 0 to 1, where 1 is fully inserted
	reactivity          float64 // 0 to 1, where 1 is critical
	criticality         float64 // 0 to 1, where 1 is critical	
	neutronFlux         float64 // in neutrons per second
	temperature         float64 // in degrees Celsius
	heatEnergyRate      float64 // in MW
	primaryLoop         *PrimaryLoop
}

func NewReactorCore(name string) *ReactorCore {
	return &ReactorCore{
		BaseComponent:       BaseComponent{Name: name},
		controlRodInsertion: 1.0, // Start with control rods fully inserted
		reactivity:          0.0,
		criticality:         0.0,
		neutronFlux:         0.0,
		temperature:         20.0, // Start at room temperature (Celsius)
		heatEnergyRate:      0.0,
	}
}

func (rc *ReactorCore) ConnectToPrimaryLoop(primaryLoop *PrimaryLoop) {
	rc.primaryLoop = primaryLoop
}

func (rc *ReactorCore) Update(env *Environment, s *Simulation) {
	// Update reactivity based on control rod insertion
	rc.reactivity = 1.0 - rc.controlRodInsertion
	rc.criticality = 1.0 - rc.controlRodInsertion

	// Update heat energy rate based on reactivity
	rc.heatEnergyRate = 3000.0 * rc.criticality // Assuming max output of 3000 MW

	// Simple temperature model (this should be more complex in reality)
	rc.temperature += (rc.heatEnergyRate / 1000.0) * 0.1          // Simplified heating
	rc.temperature = math.Max(20, math.Min(rc.temperature, 1000)) // Limit temperature range
}

func (rc *ReactorCore) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":                rc.Name,
		"controlRodInsertion": rc.controlRodInsertion,
		"criticality":         rc.criticality,
		"temperature":         rc.temperature,
		"heatEnergyRate":      rc.heatEnergyRate,
	}
}

func (rc *ReactorCore) PrintStatus() {
	fmt.Printf("Reactor Core: %s\n", rc.Name)
	fmt.Printf("\tControl Rod Insertion: %.2f\n", rc.controlRodInsertion)
	fmt.Printf("\tCriticality: %.2f\n", rc.criticality)
	fmt.Printf("\tTemperature: %.2f°C\n", rc.temperature)
	fmt.Printf("\tHeat Energy Rate: %.2f MW\n", rc.heatEnergyRate)
}

func (rc *ReactorCore) SetControlRodInsertion(insertion float64) {
	rc.controlRodInsertion = math.Max(0, math.Min(insertion, 1))
}

func (rc *ReactorCore) GetHeatEnergyRate() float64 {
	return rc.heatEnergyRate
}
