package sim

import (
	"encoding/json"
	"fmt"
	"os"
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
	verbose     bool
	stopChan    chan struct{}
	history     []map[string]interface{} // New field to store history
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
			PowerOn: true,
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

func (s *Simulation) SetVerboseLogging(verbose bool) {
	s.verbose = verbose
}

func (s *Simulation) FindPrimaryLoop() *PrimaryLoop {
	for _, component := range s.components {
		if primaryLoop, ok := component.(*PrimaryLoop); ok {
			return primaryLoop
		}
	}
	return nil
}

func (s *Simulation) FindReactorCore() *ReactorCore {
	for _, component := range s.components {
		if reactorCore, ok := component.(*ReactorCore); ok {
			return reactorCore
		}
	}
	return nil
}

func (s *Simulation) FindPressurizer() *Pressurizer {
	for _, component := range s.components {
		if pressurizer, ok := component.(*Pressurizer); ok {
			return pressurizer
		}
	}
	return nil
}

func (s *Simulation) FindSteamGenerator() *SteamGenerator {
	for _, component := range s.components {
		if steamGenerator, ok := component.(*SteamGenerator); ok {
			return steamGenerator
		}
	}
	return nil
}

func (s *Simulation) FindSecondaryLoop() *SecondaryLoop {
	for _, component := range s.components {
		if secondaryLoop, ok := component.(*SecondaryLoop); ok {
			return secondaryLoop
		}
	}
	return nil
}

func (s *Simulation) FindSteamTurbine() *SteamTurbine {
	for _, component := range s.components {
		if turbine, ok := component.(*SteamTurbine); ok {
			return turbine
		}
	}
	return nil
}

func (s *Simulation) FindCondenser() *Condenser {
	for _, component := range s.components {
		if condenser, ok := component.(*Condenser); ok {
			return condenser
		}
	}
	return nil
}

func (s *Simulation) FindGenerator() *Generator {
	for _, component := range s.components {
		if generator, ok := component.(*Generator); ok {
			return generator
		}
	}
	return nil
}

func (s *Simulation) updateEnvironment() {
	weathers := []string{"Sunny", "Cloudy", "Rainy", "Windy"}
	s.environment.Weather = weathers[s.clock.currentIter%len(weathers)]
}

func (s *Simulation) Status() map[string]interface{} {
	status := map[string]interface{}{
		"id":              s.info.ID,
		"name":            s.info.Name,
		"motto":           s.info.Motto,
		"spawned_at":      s.info.SpawnedAt,
		"simTime":         s.clock.SimTime(),
		"iterationNumber": s.clock.currentIter,
		"running":         s.running,
		"powerOn":         s.environment.PowerOn,
		"weather":         s.environment.Weather,
		"components":      make([]map[string]interface{}, 0),
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
	fmt.Printf("Power On: %t\n", s.environment.PowerOn)
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
		if s.verbose {
			fmt.Println("Whew. That was a nice run.")
			s.PrintStatus()
		}
	}()

	if s.verbose {
		fmt.Printf("Starting %d iterations\n", ticks)
	}
	for cnt = 0; cnt < ticks; cnt++ {
		select {
		case <-s.stopChan:
			if s.verbose {
				fmt.Printf("Interrupted after %d iterations\n", cnt)
			}
			return
		default:
			s.clock.Tick()
			s.updateEnvironment()
			for _, component := range s.components {
				component.Update(&s.environment, s)
			}
			s.logCurrentState() // Log the current state
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

// New method to log the current state
func (s *Simulation) logCurrentState() {
	currentStatus := s.Status()                  // Get current status
	s.history = append(s.history, currentStatus) // Append to history
	err := s.WriteHistoryToFile("simulation_history.json")
	if err != nil {
		fmt.Println("Error writing history to file:", err)
	}
}

// New method to get history
func (s *Simulation) GetHistory() []map[string]interface{} {
	return s.history
}

func (s *Simulation) WriteHistoryToFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, entry := range s.history {
		data, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		_, err = file.Write(append(data, '\n')) // Write each entry as a new line
		if err != nil {
			return err
		}
	}
	return nil
}
