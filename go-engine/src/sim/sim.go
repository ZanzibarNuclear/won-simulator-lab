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
		environment: common.Environment{
			Weather: "Sunny", // Initialize with a default weather
		},
	}
}

func (s *Simulation) AddComponent(p components.Component) {
	s.components = append(s.components, p)
}

func (s *Simulation) run(ticks int) {
	s.running = true
	defer func() { s.running = false }()

	for i := 0; i < ticks; i++ {
		s.clock.Tick()
		fmt.Printf("Iteration %d\n", s.clock.currentIter)

		s.updateEnvironment()

		for _, component := range s.components {
			component.Update(&s.environment, s.components)
		}
		s.PrintStatus()
	}
}

// Add this new method to update the environment
func (s *Simulation) updateEnvironment() {
	// This is a simple example. You might want to implement more complex weather patterns
	weathers := []string{"Sunny", "Cloudy", "Rainy", "Windy"}
	s.environment.Weather = weathers[s.clock.currentIter%len(weathers)]
}

func (s *Simulation) Status() map[string]interface{} {
	status := map[string]interface{}{
		"running":     s.running,
		"currentTime": s.clock.SimTime(),
		"weather":     s.environment.Weather,
		"components":  make([]map[string]interface{}, 0),
	}

	for _, component := range s.components {
		componentStatus := component.Status()
		status["components"] = append(status["components"].([]map[string]interface{}), componentStatus)
	}

	return status
}

func (s *Simulation) PrintStatus() {
	fmt.Println(s.Status())
	// for _, component := range s.components {
	// 	component.PrintStatus()
	// }
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
