package sim

// FIXME: perhaps this is not useful in Go. experiment with "functional options" pattern instead.

// SimulationConfigurator implements the builder pattern for configuring a Simulation
type SimulationConfigurator struct {
	simulation *Simulation
}

// NewSimulationConfigurator creates a new SimulationConfigurator
func NewSimulationConfigurator(name, motto string) *SimulationConfigurator {
	return &SimulationConfigurator{
		simulation: NewSimulation(name, motto),
	}
}

// Build returns the configured Simulation
func (c *SimulationConfigurator) Build() *Simulation {
	return c.simulation
}

// SetVerboseLogging sets the verbose logging option
func (c *SimulationConfigurator) WithVerboseLogging() *SimulationConfigurator {
	c.simulation.SetVerboseLogging(true)
	return c
}

// AddPrimaryLoop adds a PrimaryLoop to the Simulation
func (c *SimulationConfigurator) AddPrimaryLoop(name string) *SimulationConfigurator {
	c.simulation.AddComponent(NewPrimaryLoop(name))
	return c
}

// AddReactorCore adds a ReactorCore to the Simulation
func (c *SimulationConfigurator) AddReactorCore(name string) *SimulationConfigurator {
	c.simulation.AddComponent(NewReactorCore(name))
	return c
}

// AddPressurizer adds a Pressurizer to the Simulation
func (c *SimulationConfigurator) AddPressurizer(name string) *SimulationConfigurator {
	c.simulation.AddComponent(NewPressurizer(name))
	return c
}

// AddSteamGenerator adds a SteamGenerator to the Simulation
func (c *SimulationConfigurator) AddSteamGenerator(name string) *SimulationConfigurator {
	c.simulation.AddComponent(NewSteamGenerator(name))
	return c
}

// AddSecondaryLoop adds a SecondaryLoop to the Simulation
func (c *SimulationConfigurator) AddSecondaryLoop(name string) *SimulationConfigurator {
	c.simulation.AddComponent(NewSecondaryLoop(name))
	return c
}

// AddSteamTurbine adds a SteamTurbine to the Simulation
func (c *SimulationConfigurator) AddSteamTurbine(name string) *SimulationConfigurator {
	c.simulation.AddComponent(NewSteamTurbine(name))
	return c
}

// AddCondenser adds a Condenser to the Simulation
func (c *SimulationConfigurator) AddCondenser(name string) *SimulationConfigurator {
	c.simulation.AddComponent(NewCondenser(name))
	return c
}

// AddGenerator adds a Generator to the Simulation
func (c *SimulationConfigurator) AddGenerator(name string) *SimulationConfigurator {
	c.simulation.AddComponent(NewGenerator(name))
	return c
}
