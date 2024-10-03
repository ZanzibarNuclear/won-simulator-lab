package pwr

import (
	"fmt"
	"math"

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

	sl := NewSecondaryLoop("SL1", "The secondary loop relies primarily on convection to circulate forms of water through the steam generators, steam turbine, and condenser. Feedwater pumps supply the steam generator with replacement water.")
	s.AddComponent(sl)

	rc := NewReactorCore("RC1", "The reactor core is the heart of the PWR, where the nuclear reaction takes place, creating massive amounts of heat. The control rods are encapsulated in this component.")
	s.AddComponent(rc)

	pr := NewPressurizer("PR1", "The pressurizer creates a bubble of steam to pressurize the primary loop. This keeps high-temperature water in liquid state as it is pumped through the reactor core to the steam generators.")
	s.AddComponent(pr)

	sg := NewSteamGenerator("SG1", "The steam generator enables the transfer of heat from the hot leg of the primary loop to the water under lower pressure in the secondary loop.", pl, sl)
	s.AddComponent(sg)

	st := NewSteamTurbine("ST1", "The steam turbine converts heat energy to mechanical energy as steam is forced through it, pushing on the turbine blades and causing them to rotate. The turbine then drives the generator.", sl)
	s.AddComponent(st)

	g := NewGenerator("G1", "The generator converts mechanical energy to electrical energy via electromagnetic induction.", st)
	s.AddComponent(g)

	c := NewCondenser("C1", "The condenser extracts heat from the steam as it leaves the steam turbine, turning it back into water that flows back to the feedwater pumps.", st)
	s.AddComponent(c)
}

func (s *PwrSim) PrintStatus() {
	println("\n--- Like we always say: ", s.motto, " ---")
	s.Simulator.PrintStatus()
	println("----------------------------------------")
	println("----------------------------------------\n\n")
}

func (s *PwrSim) ProcessEvent(event *simworks.Event) {
	switch event.Code {
	case Event_pr_reliefValveVent:
		fmt.Println("The pressurizer relief valve was triggered and is venting.")
		event.SetComplete()
	case Event_sl_emergencyMssvVent:
		fmt.Println("The emergency MSSV was released and is venting.")
		event.SetComplete()
	default:
		fmt.Println("Event code not handled:", event.Code)
	}

}

func CalcLinearIncrease(currentValue, targetValue, stepSize float64) float64 {
	if currentValue < targetValue {
		return currentValue + math.Min(targetValue-currentValue, stepSize)
	}
	return targetValue
}

func CalcLinearDecrease(currentValue, targetValue, stepSize float64) float64 {
	if currentValue > targetValue {
		return currentValue - math.Min(currentValue-targetValue, stepSize)
	}
	return targetValue
}
