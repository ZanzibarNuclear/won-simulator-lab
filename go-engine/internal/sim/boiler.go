package sim

import (
	"fmt"
)

// Example part implementation
type Boiler struct {
	BaseComponent
	running           bool
	targetTemperature int
	temperature       int
	heatEnergy        int
	fuel              FuelConsumption
}

type FuelConsumption struct {
	consumed       int
	rateOfIncrease int
	rate           int
}

func NewBoiler(name string) *Boiler {
	return &Boiler{
		BaseComponent:     BaseComponent{Name: name},
		running:           false,
		temperature:       ROOM_TEMPERATURE,
		targetTemperature: 600,
		heatEnergy:        0,
		fuel: FuelConsumption{
			consumed:       0,
			rateOfIncrease: 0,
			rate:           0,
		},
	}
}

func (p *Boiler) Update(env *Environment, s *Simulation) {

	turbine := s.FindTurbine()
	if turbine == nil {
		fmt.Println("No turbine found")
		return
	}

	// simulating start-up to steady state
	// increase rate of consumption until turbine reaches its limit
	if !turbine.MaxedOut() && p.temperature < p.targetTemperature {
		p.fuel.rate += p.fuel.rateOfIncrease
		// FIXME: replace with a model that reflects temperature change due to net heat input
		p.temperature += 20
	} else {
		p.fuel.rateOfIncrease = 0
	}
	p.fuel.consumed += p.fuel.rate
	p.heatEnergy += p.fuel.rate*10 ^ 7

	// FIXME: ultimately, want boiler to be controlled by operator (which could be AI)
}

func (p *Boiler) PrintStatus() {
	fmt.Printf("Boiler: %s\n", p.Name)
	fmt.Printf("\tRunning: %t\n", p.running)
	fmt.Printf("\tTemperature: %d\n", p.temperature)
	fmt.Printf("\tHeat power: %d\n", p.heatEnergy)
	fmt.Println("\tFuel consumption:")
	fmt.Printf("\t\tCurrent rate (per iteration): %d\n", p.fuel.rate)
	fmt.Printf("\t\tRate of increase: %d\n", p.fuel.rateOfIncrease)
	fmt.Printf("\t\tCumulative: %d\n", p.fuel.consumed)
	fmt.Println()
}

func (p *Boiler) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":        p.Name,
		"running":     p.running,
		"temperature": p.temperature,
		"heatEnergy":  p.heatEnergy,
		"fuelConsumption": map[string]int{
			"rate":           p.fuel.rate,
			"rateOfIncrease": p.fuel.rateOfIncrease,
			"consumed":       p.fuel.consumed,
		},
	}
}

func (p *Boiler) Running() bool {
	return p.running
}

func (p *Boiler) TurnOn() {
	p.running = true
	p.fuel.rateOfIncrease = 1
}

func (p *Boiler) TurnOff() {
	p.running = false
	p.fuel.rateOfIncrease = 0
}

func (p *Boiler) TurnUp() {
	p.fuel.rateOfIncrease += 1
}

func (p *Boiler) TurnDown() {
	if p.fuel.rateOfIncrease > 0 {
		p.fuel.rateOfIncrease -= 1
	}
}

func (p *Boiler) HoldSteady() {
	p.fuel.rateOfIncrease += 0
}
