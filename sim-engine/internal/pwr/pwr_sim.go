package pwr

import (
	"worldofnuclear.com/internal/simworks"
)

type PwrSim struct {
	simworks.Simulator
	motto string
}

func NewPwrSim(name, motto string) *PwrSim {
	return &PwrSim{
		Simulator: *simworks.NewSimulator(name, "Emulate a Pressurized Water Reactor. This is for demonstration and entertainment purposes and is not suitable for running a real PWR."),
		motto:     motto,
	}
}

func (s *PwrSim) SetupStandardComponents() {
	pl := NewPrimaryLoop("PL1", "The primary loop circulates coolant through the reactor core, carrying boron for moderation and heat from the core to the steam generators.")
	s.AddComponent(pl)
}

func (s *PwrSim) PrintStatus() {
	s.Simulator.PrintStatus()
	println("\n--- Like we always say: ", s.motto, " ---")
}
