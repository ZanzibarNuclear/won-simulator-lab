package main

import (
	"worldofnuclear.com/internal/pwr"
)

func main() {
	sim := pwr.NewPwrSim("PWR Sim", "Some like it hot.")
	sim.SetEventHandler(sim)

	sim.SetupStandardComponents()
	sim.Step()

	sim.QueueEvent(pwr.NewEvent_PumpSwitch(true))
	sim.RunForABit(0, 0, 1, 0)
	// sim.PrintStatus()

	sim.QueueEvent(pwr.NewEvent_BoronConcentration(300.0))
	sim.RunForABit(0, 0, 25, 0)
	// sim.PrintStatus()

	sim.QueueEvent(pwr.NewEvent_TargetPressure(15.5))
	sim.QueueEvent(pwr.NewEvent_HeaterPower(true))
	sim.RunForABit(0, 0, 0, 10)
	// sim.PrintStatus()

	sim.RunForABit(0, 3, 0, 0)
	sim.PrintStatus()
}

// func sampleSim() {
// 	sim := simworks.NewSimulator("Simulator Works", "The inner workings to support simulations")
// 	pl := pwr.NewPrimaryLoop("TestLoop-Pump", "The is a test.")
// 	sim.AddComponent(pl)
// 	sim.RunForABit(0, 0, 1, 30)

// 	pumpOn := simworks.NewImmediateEventBool(pwr.Event_pl_pumpSwitch, true)
// 	sim.QueueEvent(pumpOn)
// 	sim.Step()

// 	if !pl.PumpOn() {
// 		fmt.Println("Pump is not on. Boo hoo.")
// 	}

// 	targetBoron := 300.0
// 	addBoron := simworks.NewAdjustmentEvent(pwr.Event_pl_boronConcentration, targetBoron)
// 	sim.QueueEvent(addBoron)

// 	sim.Run(1) // should pick up event; will need some time to complete

// 	for i := 0; i < 4 && !addBoron.IsComplete(); i++ {
// 		sim.RunForABit(0, 0, 30, 0)
// 		fmt.Printf("%s Boron concentration: %f\n", sim.CurrentMoment(), pl.BoronConcentration())
// 	}

// 	if pl.BoronConcentration() != targetBoron {
// 		fmt.Printf("Boron concentration did not reach target: %f. Boo hoo: %f.\n", targetBoron, pl.BoronConcentration())
// 	}
// }
