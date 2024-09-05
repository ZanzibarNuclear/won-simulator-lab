package sim

import (
	"fmt"
	"time"
	"won/sim-lab/go-engine/common"
	"won/sim-lab/go-engine/components"
)

type Clock struct {
	startedAt   time.Time
	currentIter int
}

func (c *Clock) SimTime() time.Time {
	return c.startedAt.Add(time.Duration(c.currentIter) * time.Minute)
}

func (c *Clock) Tick() {
	c.currentIter++
}

type Simulation struct {
	components  []components.Component
	clock       Clock
	environment common.Environment
	running     bool
}

func NewSimulation() *Simulation {
	return &Simulation{
		components: make([]components.Component, 0),
		clock: Clock{
			startedAt:   time.Date(2000, 1, 1, 8, 0, 0, 0, time.FixedZone("EST", -5*60*60)),
			currentIter: 0,
		},
	}
}

func (s *Simulation) AddPart(p components.Component) {
	s.components = append(s.components, p)
}

func (s *Simulation) run(ticks int) {
	s.running = true
	defer func() { s.running = false }()

	for i := 0; i < ticks; i++ {
		s.clock.Tick()
		fmt.Printf("Iteration %d\n", s.clock.currentIter)

		for _, component := range s.components {
			component.Update(&s.environment, s.components)
		}
		s.PrintStatus()
	}
}

func (s *Simulation) GetStatus() string {
	status := fmt.Sprintf("\n=====\nSimulation Status at %s\n", s.clock.SimTime().Format("2006-01-02 15:04"))
	status += fmt.Sprintf("\tRunning: %t\n", s.running)
	status += "\tEnvironment:\n"
	status += fmt.Sprintf("\t\tWeather: %s\n", s.environment.Weather)
	return status
}

func (s *Simulation) PrintStatus() {
	fmt.Println(s.GetStatus())
	for _, component := range s.components {
		component.PrintStatus()
	}
}

func (s *Simulation) IsRunning() bool {
	return s.running
}

func (s *Simulation) CurrentTime() time.Time {
	return s.clock.SimTime()
}

func (s *Simulation) Advance(iterations int) {
	if !s.running {
		s.run(iterations)
	}
}

func (s *Simulation) Start() {
	if !s.running {
		go s.run(common.YEAR_OF_MINUTES) // Run for a year by default
	}
}

func (s *Simulation) Stop() {
	// TODO: Implement stopping the simulation
}
