package main

import (
	"fmt"
	"won/sim-lab/go-engine/common"
	"won/sim-lab/go-engine/components"
)

type Simulation struct {
	parts       []components.Component
	environment common.Environment
	iterations  int
}

func NewSimulation(iterations int) *Simulation {
	return &Simulation{
		parts:      make([]components.Component, 0),
		iterations: iterations,
	}
}

func (s *Simulation) AddPart(p components.Component) {
	s.parts = append(s.parts, p)
}

func (s *Simulation) Run() {
	for i := 0; i < s.iterations; i++ {
		fmt.Printf("Iteration %d\n", i+1)
		for _, part := range s.parts {
			part.Update(&s.environment, s.parts)
		}
		fmt.Println()
		// Update environment if needed
	}
}

func main() {
	hour := 60
	day := hour * 24
	week := day * 7
	year := week * 52
	sim := NewSimulation(year) // minutes per day

	// Add parts to the simulation
	// sim.AddPart(parts.NewSomePart(...))
	boiler := components.NewBoiler()
	steamTurbine := components.NewTurbine()

	sim.AddPart(boiler)
	sim.AddPart(steamTurbine)

	sim.Run()
}
