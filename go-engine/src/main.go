package main

import (
	"log"
	"won/sim-lab/go-engine/common"
	"won/sim-lab/go-engine/components"
	"won/sim-lab/go-engine/server"
	"won/sim-lab/go-engine/sim"
)

func main() {
	sim := sim.NewSimulation(common.HOUR_OF_MINUTES) // one day of one-minute iterations

	// Add parts to the simulation
	// sim.AddPart(parts.NewSomePart(...))
	boiler := components.NewBoiler()
	boiler.TurnOn()
	steamTurbine := components.NewTurbine()

	sim.AddPart(boiler)
	sim.AddPart(steamTurbine)

	sim.Run()

	svr := server.NewServer(sim)
	log.Fatal(svr.Start())
}
