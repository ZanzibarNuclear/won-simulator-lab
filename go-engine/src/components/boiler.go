package components

import (
	"fmt"
	"won/sim-lab/go-engine/common"
)

// Example part implementation
type Boiler struct {
	running           bool
	targetTemperature int
	temperature       int
	heatEnergy        int
	fuelConsumption   FuelConsumption
}

type FuelConsumption struct {
	consumed       int
	rateOfIncrease int
	rate           int
}

func NewBoiler() *Boiler {
	return &Boiler{
		running:           false,
		temperature:       common.ROOM_TEMPERATURE,
		targetTemperature: 600,
		heatEnergy:        0,
		fuelConsumption: FuelConsumption{
			consumed:       0,
			rateOfIncrease: 0,
			rate:           0,
		},
	}
}

func (p *Boiler) Update(env *common.Environment, otherComponents []Component) {

	turbine := FindTurbine(otherComponents)
	if turbine == nil {
		fmt.Println("No turbine found")
		return
	}

	// simulating start-up to steady state
	// increase rate of consumption until turbine reaches its limit
	if !turbine.MaxedOut() && p.temperature < p.targetTemperature {
		p.fuelConsumption.rate += p.fuelConsumption.rateOfIncrease
		// FIXME: replace with a model that reflects temperature change due to net heat input
		p.temperature += 20
	} else {
		p.fuelConsumption.rateOfIncrease = 0
	}
	p.fuelConsumption.consumed += p.fuelConsumption.rate
	p.heatEnergy += p.fuelConsumption.rate*10 ^ 7

	// FIXME: ultimately, want boiler to be controlled by operator (which could be AI)
}

func (p *Boiler) PrintStatus() {
	fmt.Println("Boiler status:")
	fmt.Printf("\tRunning: %t\n", p.running)
	fmt.Printf("\tTemperature: %d\n", p.temperature)
	fmt.Printf("\tHeat power: %d\n", p.heatEnergy)
	fmt.Println("\tFuel consumption:")
	fmt.Printf("\t\tCurrent rate (per iteration): %d\n", p.fuelConsumption.rate)
	fmt.Printf("\t\tRate of increase: %d\n", p.fuelConsumption.rateOfIncrease)
	fmt.Printf("\t\tCumulative: %d\n", p.fuelConsumption.consumed)
	fmt.Println()
}

func (p *Boiler) Status() map[string]interface{} {
	return map[string]interface{}{
		"running":     p.running,
		"temperature": p.temperature,
		"heatEnergy":  p.heatEnergy,
		"fuelConsumption": map[string]int{
			"rate":           p.fuelConsumption.rate,
			"rateOfIncrease": p.fuelConsumption.rateOfIncrease,
			"consumed":       p.fuelConsumption.consumed,
		},
	}
}

func FindBoiler(components []Component) *Boiler {
	for _, component := range components {
		if boiler, ok := component.(*Boiler); ok {
			return boiler
		}
	}
	return nil
}

func (p *Boiler) Running() bool {
	return p.running
}

func (p *Boiler) TurnOn() {
	p.running = true
	p.fuelConsumption.rateOfIncrease = 1
}

func (p *Boiler) TurnOff() {
	p.running = false
	p.fuelConsumption.rateOfIncrease = 0
}

func (p *Boiler) TurnUp() {
	p.fuelConsumption.rateOfIncrease += 1
}

func (p *Boiler) TurnDown() {
	if p.fuelConsumption.rateOfIncrease > 0 {
		p.fuelConsumption.rateOfIncrease -= 1
	}
}

func (p *Boiler) HoldSteady() {
	p.fuelConsumption.rateOfIncrease += 0
}
