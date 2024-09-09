package sim

import (
	"fmt"
	"math"
)

type Turbine struct {
	BaseComponent
	maxRpm int
	rpm    int
}

func NewTurbine(name string) *Turbine {
	return &Turbine{
		BaseComponent: BaseComponent{Name: name},
		maxRpm:        TURBINE_MAX_RPM,
		rpm:           0,
	}
}

func (t *Turbine) Update(env *Environment, s *Simulation) {
	boiler := s.FindBoiler()
	if boiler == nil {
		fmt.Println("No boiler found")
		return
	}
	// how fast based on fuel consumption
	t.rpm = int(math.Min(float64(t.maxRpm), float64(boiler.fuel.rate*500)))
}

func (t *Turbine) Status() map[string]interface{} {
	return map[string]interface{}{
		"name": t.Name,
		"rpm": t.rpm,
	}
}

func (t *Turbine) PrintStatus() {
	fmt.Printf("Turbine: %s\n", t.Name)
	fmt.Printf("\tSpinning at %d RPMs\n", t.rpm)
	if t.MaxedOut() {
		fmt.Printf("\tMaxed out!!\n")
	}
	fmt.Println()
}

func (t *Turbine) MaxedOut() bool {
	return t.rpm >= t.maxRpm
}

func (t *Turbine) Rpm() int {
	return t.rpm
}
