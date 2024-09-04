package components

import (
	"fmt"
	"won/sim-lab/go-engine/common"
)

// Example part implementation
type Boiler struct {
	fuelConsumption FuelConsumption
}

type FuelConsumption struct {
	consumed       int
	rateOfIncrease int
	rate           int
}

func NewBoiler() *Boiler {
	return &Boiler{
		fuelConsumption: FuelConsumption{
			consumed:       0,
			rateOfIncrease: 1,
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
	if !turbine.IsMaxedOut() {
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
