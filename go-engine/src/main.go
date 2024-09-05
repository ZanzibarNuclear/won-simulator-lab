package main

import (
	// "won/sim-lab/go-engine/common"
	"won/sim-lab/go-engine/components"
	"won/sim-lab/go-engine/sim"
)

func main() {
	sim := sim.NewSimulation() // one day of one-minute iterations

	// Add parts to the simulation
	boiler := components.NewBoiler()
	boiler.TurnOn()

	steamTurbine := components.NewTurbine()

	sim.AddPart(boiler)
	sim.AddPart(steamTurbine)

	sim.Advance(5)

	// svr := server.NewServer(sim)
	// log.Fatal(svr.Start())
}
