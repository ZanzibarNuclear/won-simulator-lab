package parts

import (
	"fmt"
	"math"
	"won/sim-lab/go-engine/common"
)

type Part interface {
	Update(env *common.Environment, otherParts []Part)
}

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

func (p *Boiler) Update(env *common.Environment, otherParts []Part) {

	turbine := FindTurbine(otherParts)
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

func FindBoiler(parts []Part) *Boiler {
	for _, part := range parts {
		if boiler, ok := part.(*Boiler); ok {
			return boiler
		}
	}
	return nil
}

type Turbine struct {
	maxRpm int
	rpm    int
}

func NewTurbine() *Turbine {
	return &Turbine{
		maxRpm: 3600,
		rpm:    0,
	}
}

func (t *Turbine) Update(env *common.Environment, otherParts []Part) {
	boiler := FindBoiler(otherParts)
	if boiler == nil {
		fmt.Println("No boiler found")
		return
	}
	// how fast based on fuel consumption
	t.rpm = int(math.Min(float64(t.maxRpm), float64(boiler.fuelConsumption.rate*500)))

	info := fmt.Sprintf("Turbine spinning at %d RPMs.", t.rpm)
	fmt.Println(info)

	if t.IsMaxedOut() {
		fmt.Println("\tTo the max!!")
	}
}

func FindTurbine(parts []Part) *Turbine {
	for _, part := range parts {
		if turbine, ok := part.(*Turbine); ok {
			return turbine
		}
	}
	return nil
}

func (t *Turbine) IsMaxedOut() bool {
	return t.rpm >= t.maxRpm
}
