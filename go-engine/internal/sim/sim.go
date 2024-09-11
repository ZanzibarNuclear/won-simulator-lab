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

	stopChan chan struct{}
}

func NewSimulation(name string, motto string) *Simulation {
	return &Simulation{
		info: SimInfo{
			ID:        fmt.Sprintf("sim-%s", generateRandomID(8)),
			Name:      name,
			Motto:     motto,
			SpawnedAt: time.Now(),
		},
		clock: Clock{
			startedAt:   time.Date(2000, 1, 1, 8, 0, 0, 0, time.FixedZone("EST", -5*60*60)),
			currentIter: 0,
		},
		environment: Environment{
			Weather: "Sunny", // Initialize with a default weather
		},
		components: make([]Component, 0),
		stopChan:   make(chan struct{}),
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
		"id":         s.info.ID,
		"name":       s.info.Name,
		"motto":      s.info.Motto,
		"spawned_at": s.info.SpawnedAt,
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
	fmt.Println("----------------------------------------")
	fmt.Printf("Name: %s\n", s.info.Name)
	fmt.Printf("ID: %s\n", s.info.ID)
	fmt.Printf("Motto: %s\n", s.info.Motto)
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

func (s *Simulation) Run(ticks int) {
	// one running process at a time
	if s.running {
		return
	}
	var cnt int
	s.running = true

	defer func() {
		s.running = false
		fmt.Println("Whew. That was a nice run.")
		s.PrintStatus()
	}()

	fmt.Printf("Starting %d iterations\n", ticks)
	for cnt = 0; cnt < ticks; cnt++ {
		select {
		case <-s.stopChan:
			fmt.Printf("Interrupted after %d iterations\n", cnt)
			return
		default:
			s.clock.Tick()
			s.updateEnvironment()
			for _, component := range s.components {
				component.Update(&s.environment, s)
			}
		}
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
		s.stopChan = make(chan struct{})
		go s.Run(iterations)
	}
}

func (s *Simulation) AdvanceOneYear() {
	s.Advance(YEAR_OF_MINUTES)
}

func (s *Simulation) Stop() {
	if s.running {
		close(s.stopChan)
	}
}
