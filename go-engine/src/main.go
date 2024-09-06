package main

import (
	"won/sim-lab/go-engine/components"
	"won/sim-lab/go-engine/server"
	"won/sim-lab/go-engine/sim"
	"log"
)

func main() {
	sim := sim.NewSimulation() // one day of one-minute iterations

	// Add parts to the simulation
	boiler := components.NewBoiler()
	boiler.TurnOn()

	steamTurbine := components.NewTurbine()

	sim.AddPart(boiler)
	sim.AddPart(steamTurbine)

	// sim.Advance(common.WEEK_OF_MINUTES)

	svr := server.NewServer(sim)
	log.Fatal(svr.Start())
}
