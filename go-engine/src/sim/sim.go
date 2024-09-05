package sim

import (
	"fmt"
	"won/sim-lab/go-engine/common"
	"won/sim-lab/go-engine/components"
)

type Simulation struct {
	components  []components.Component
	environment common.Environment
	iterations  int
	running     bool
}

func NewSimulation(iterations int) *Simulation {
	return &Simulation{
		components: make([]components.Component, 0),
		iterations: iterations,
	}
}

func (s *Simulation) AddPart(p components.Component) {
	s.components = append(s.components, p)
}

func (s *Simulation) Run() {
	for i := 0; i < s.iterations; i++ {
		fmt.Printf("Iteration %d\n", i+1)
		for _, component := range s.components {
			component.Update(&s.environment, s.components)
		}
		fmt.Println()
		// Update environment if needed
	}
}

func (s *Simulation) IsRunning() bool {
	return s.running
}

func (s *Simulation) Start() {
	s.running = true
}

func (s *Simulation) Stop() {
	s.running = false
}