package main

import (
	"log"
	"won/sim-lab/go-engine/internal/components"
	"won/sim-lab/go-engine/internal/server"
	"won/sim-lab/go-engine/internal/sim"
)

func main() {
	sim := sim.NewSimulation() // one day of one-minute iterations

	// Add parts to the simulation
	boiler := components.NewBoiler("Billy Boyle")
	boiler.TurnOn()

	steamTurbine := components.NewTurbine("Tilly Turner")

	sim.AddComponent(boiler)
	sim.AddComponent(steamTurbine)

	// sim.Advance(common.WEEK_OF_MINUTES)

	svr := server.NewServer(sim)
	log.Fatal(svr.Start())
}
