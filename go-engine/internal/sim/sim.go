package sim

import (
	"fmt"
	"time"
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
	info        SimInfo
	components  []Component
	clock       Clock
	environment Environment
	running     bool
}

func NewSimulation(name string, motto string) *Simulation {
	return &Simulation{
		info: SimInfo{
			ID:        fmt.Sprintf("sim-%s", generateRandomID(8)),
			Name:      name,
			Motto:     motto,
			SpawnedAt: time.Now(),
		},
		components: make([]Component, 0),
		clock: Clock{
			startedAt:   time.Date(2000, 1, 1, 8, 0, 0, 0, time.FixedZone("EST", -5*60*60)),
			currentIter: 0,
		},
		environment: Environment{
			Weather: "Sunny", // Initialize with a default weather
		},
	}
}

type SimInfo struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Motto     string    `json:"motto"`
	SpawnedAt time.Time `json:"spawned_at"`
}

func (s *Simulation) AddComponent(p Component) {
	s.components = append(s.components, p)
}

func (s *Simulation) run(ticks int) {
	s.running = true
	defer func() { s.running = false }()

	for i := 0; i < ticks; i++ {
		s.clock.Tick()
		fmt.Printf("Starting iteration %d\n", s.clock.currentIter)

		s.updateEnvironment()

		for _, component := range s.components {
			component.Update(&s.environment, s)
		}
		s.PrintStatus()
	}
}

func (s *Simulation) ID() string {
	return s.info.ID
}

func (s *Simulation) Info() SimInfo {
	return s.info
}

func (s *Simulation) Components() []Component {
	return s.components
}

func (s *Simulation) FindBoiler() *Boiler {
	for _, component := range s.components {
		if boiler, ok := component.(*Boiler); ok {
			return boiler
		}
	}
	return nil
}

func (s *Simulation) FindTurbine() *Turbine {
	for _, component := range s.components {
		if turbine, ok := component.(*Turbine); ok {
			return turbine
		}
	}
	return nil
}

// Add this new method to update the environment
func (s *Simulation) updateEnvironment() {
	// This is a simple example. You might want to implement more complex weather patterns
	weathers := []string{"Sunny", "Cloudy", "Rainy", "Windy"}
	s.environment.Weather = weathers[s.clock.currentIter%len(weathers)]
}

func (s *Simulation) Status() map[string]interface{} {
	status := map[string]interface{}{
		"running":    s.running,
		"simTime":    s.clock.SimTime(),
		"weather":    s.environment.Weather,
		"components": make([]map[string]interface{}, 0),
	}
	for _, component := range s.components {
		componentStatus := component.Status()
		status["components"] = append(status["components"].([]map[string]interface{}), componentStatus)
	}
	return status
}

func (s *Simulation) PrintStatus() {
	fmt.Printf("Sim Time: %s\n", s.clock.SimTime())
	fmt.Printf("Started at: %s\n", s.clock.startedAt)
	fmt.Printf("Is running: %t\n", s.running)
	fmt.Printf("Last iteration %d\n", s.clock.currentIter)
	fmt.Printf("Weather: %s\n\n", s.environment.Weather)
	for _, component := range s.components {
		component.PrintStatus()
	}
	fmt.Println("----------------------------------------")
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
		go s.run(YEAR_OF_MINUTES) // Run for a year by default
	}
}

func (s *Simulation) Stop() {
	// TODO: Implement stopping the simulation
}
