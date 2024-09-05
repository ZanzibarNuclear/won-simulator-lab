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
		maxRpm: common.TURBINE_MAX_RPM,
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
}

func (t *Turbine) PrintStatus() {
	fmt.Println("Turbine status:")
	fmt.Printf("\tSpinning at %d RPMs\n", t.rpm)
	if t.MaxedOut() {
		fmt.Printf("\tMaxed out!!\n")
	}
	fmt.Println()
}

func FindTurbine(components []Component) *Turbine {
	for _, component := range components {
		if turbine, ok := component.(*Turbine); ok {
			return turbine
		}
	}
	return nil
}

func (t *Turbine) MaxedOut() bool {
	return t.rpm >= t.maxRpm
}

func (t *Turbine) Rpm() int {
	return t.rpm
}
