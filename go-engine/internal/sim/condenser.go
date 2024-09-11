package sim

import (
	"fmt"
)

type Condenser struct {
	BaseComponent
	entryTemperature float64 // in Celsius
	exitTemperature  float64 // in Celsius
	heatTransferRate float64 // in Watts
}

func NewCondenser(name string) *Condenser {
	return &Condenser{
		BaseComponent:    BaseComponent{Name: name},
		entryTemperature: 100.0, // Initial values, can be adjusted as needed
		exitTemperature:  40.0,
		heatTransferRate: 1000000.0, // 1 MW, for example
	}
}

func (c *Condenser) Update(env *Environment, s *Simulation) {
	// Simplified update logic
	// In a real scenario, this would involve complex thermodynamics calculations
	turbine := s.FindSteamTurbine()
	if turbine == nil {
		fmt.Println("No turbine found")
		return
	}

	// Adjust entry temperature based on turbine exhaust (simplified)
	c.entryTemperature = 100.0 - (1000.0-float64(turbine.rpm))*0.1

	// Calculate exit temperature (simplified)
	c.exitTemperature = c.entryTemperature - (c.heatTransferRate * 0.00001)

	// Ensure temperatures don't go below ambient
	if c.exitTemperature < float64(env.AmbientTemperature) {
		c.exitTemperature = float64(env.AmbientTemperature)
	}
}

func (c *Condenser) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":             c.Name,
		"entryTemperature": c.entryTemperature,
		"exitTemperature":  c.exitTemperature,
		"heatTransferRate": c.heatTransferRate,
	}
}

func (c *Condenser) PrintStatus() {
	fmt.Printf("Condenser: %s\n", c.Name)
	fmt.Printf("\tEntry Temperature: %.2f °C\n", c.entryTemperature)
	fmt.Printf("\tExit Temperature: %.2f °C\n", c.exitTemperature)
	fmt.Printf("\tHeat Transfer Rate: %.2f W\n", c.heatTransferRate)
}
