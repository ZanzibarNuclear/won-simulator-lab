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
	pl := NewPrimaryLoop("Primary Loop", "The primary loop is the loop that circulates the coolant through the reactor core, transferring heat from the core to the steam generators.")
	s.AddComponent(pl)
}
