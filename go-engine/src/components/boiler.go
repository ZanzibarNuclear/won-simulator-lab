package components

import (
	"fmt"
	"won/sim-lab/go-engine/common"
)

// Example part implementation
type Boiler struct {
	running         bool
	temperature     int
	heatPower       int
	fuelConsumption FuelConsumption
}

type FuelConsumption struct {
	consumed       int
	rateOfIncrease int
	rate           int
}

func NewBoiler() *Boiler {
	return &Boiler{
		running:     false,
		temperature: common.ROOM_TEMPERATURE,
		heatPower:   0,
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

	// increase rate of consumption until turbine reaches its limit
	if !turbine.MaxedOut() {
		p.fuelConsumption.rate += p.fuelConsumption.rateOfIncrease
	}

	p.fuelConsumption.consumed += p.fuelConsumption.rate

	info := fmt.Sprintf("Boiler fuel consumed\n\tthis iteration:  %d\n\trunning total: %d", p.fuelConsumption.rate, p.fuelConsumption.consumed)

	fmt.Println(info)
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
	p.running = true
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
