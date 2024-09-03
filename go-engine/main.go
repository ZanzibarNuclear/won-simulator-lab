package main

import (
	"fmt"
	"won/sim-lab/go-engine/common"
	"won/sim-lab/go-engine/parts"
)

type Simulation struct {
	parts       []parts.Part
	environment common.Environment
	iterations  int
}

func NewSimulation(iterations int) *Simulation {
	return &Simulation{
		parts:      make([]parts.Part, 0),
		iterations: iterations,
	}
}

func (s *Simulation) AddPart(p parts.Part) {
	s.parts = append(s.parts, p)
}

func (s *Simulation) Run() {
	for i := 0; i < s.iterations; i++ {
		fmt.Printf("Iteration %d\n", i+1)
		for _, part := range s.parts {
			part.Update(&s.environment, s.parts)
		}
		// Update environment if needed
	}
}

func main() {
	sim := NewSimulation(10) // 100 iterations

	// Add parts to the simulation
	// sim.AddPart(parts.NewSomePart(...))
	boiler := parts.NewBoiler()
	steamTurbine := parts.NewTurbine()

	sim.AddPart(boiler)
	sim.AddPart(steamTurbine)

	sim.Run()
}