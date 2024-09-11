package sim

import (
	"fmt"
	"math"
)

type ReactorCore struct {
	BaseComponent
	temperature    float64
	maxTemperature float64
	fuelLevel      float64
	maxFuelLevel   float64
}

func NewReactorCore(name string) *ReactorCore {
	return &ReactorCore{
		BaseComponent:  BaseComponent{Name: name},
		temperature:    20.0,   // Starting at room temperature
		maxTemperature: 1000.0, // Maximum safe temperature in Celsius
		fuelLevel:      100.0,  // Starting with full fuel
		maxFuelLevel:   100.0,
	}
}

func (r *ReactorCore) Update(env *Environment, s *Simulation) {
	// Simulate temperature increase based on fuel consumption
	fuelConsumptionRate := 0.1 // Adjust as needed
	r.fuelLevel = math.Max(0, r.fuelLevel-fuelConsumptionRate)

	// Temperature increases as fuel is consumed, but starts cooling if fuel is depleted
	if r.fuelLevel > 0 {
		r.temperature = math.Min(r.maxTemperature, r.temperature+5)
	} else {
		r.temperature = math.Max(20, r.temperature-1)
	}
}

func (r *ReactorCore) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":        r.Name,
		"temperature": r.temperature,
		"fuelLevel":   r.fuelLevel,
	}
}

func (r *ReactorCore) PrintStatus() {
	fmt.Printf("Reactor Core: %s\n", r.Name)
	fmt.Printf("\tTemperature: %.2f °C\n", r.temperature)
	fmt.Printf("\tFuel Level: %.2f%%\n", r.fuelLevel)
	if r.IsOverheating() {
		fmt.Printf("\tWARNING: Reactor is overheating!\n")
	}
	if r.IsLowOnFuel() {
		fmt.Printf("\tWARNING: Reactor is low on fuel!\n")
	}
	fmt.Println()
}

func (r *ReactorCore) IsOverheating() bool {
	return r.temperature >= r.maxTemperature*0.9
}

func (r *ReactorCore) IsLowOnFuel() bool {
	return r.fuelLevel <= r.maxFuelLevel*0.1
}
