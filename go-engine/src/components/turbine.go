package components

import (
	"fmt"
	"math"
	"won/sim-lab/go-engine/common"
)

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

func (t *Turbine) Update(env *common.Environment, otherComponents []Component) {
	boiler := FindBoiler(otherComponents)
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

func FindTurbine(components []Component) *Turbine {
	for _, component := range components {
		if turbine, ok := component.(*Turbine); ok {
			return turbine
		}
	}
	return nil
}

func (t *Turbine) IsMaxedOut() bool {
	return t.rpm >= t.maxRpm
}
