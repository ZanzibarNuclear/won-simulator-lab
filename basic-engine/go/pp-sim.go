package main

import (
	"fmt"
	"math/rand"
	"time"
)

type FuelConsumption struct {
	rateOfConsumption int
	consumed          int
	generatorSpeed    *GeneratorSpeed
}

func (fc *FuelConsumption) run(done chan bool) {
	for {
		select {
		case <-done:
			return
		default:
			if !fc.generatorSpeed.isRunningAtMax() {
				if rand.Float32() < 0.5 { // 50% chance to increase
					fc.rateOfConsumption++
				}
			}
			fc.consumed += fc.rateOfConsumption
			fmt.Printf("Updating fuel consumption at %d (%d units per tick)\n", fc.consumed, fc.rateOfConsumption)
			time.Sleep(time.Second)
		}
	}
}

type GeneratorSpeed struct {
	speed            int
	maxSpeed         int
	fuelConsumption  *FuelConsumption
}

func (gs *GeneratorSpeed) run(done chan bool) {
	for {
		select {
		case <-done:
			return
		default:
			targetSpeed := gs.fuelConsumption.rateOfConsumption
			if targetSpeed > gs.maxSpeed {
				targetSpeed = gs.maxSpeed
			}
			if gs.speed < targetSpeed {
				gs.speed++
			} else if gs.speed > targetSpeed {
				gs.speed--
			}
			fmt.Printf("Updating generator speed at %d\n", gs.speed)
			time.Sleep(time.Second)
		}
	}
}

func (gs *GeneratorSpeed) isRunningAtMax() bool {
	return gs.speed >= gs.maxSpeed
}

type ElectricityDemand struct {
	demand int
}

func (ed *ElectricityDemand) run(done chan bool) {
	for {
		select {
		case <-done:
			return
		default:
			increase := rand.Intn(11) - 5 // Random number between -5 and 5
			ed.demand += increase
			if ed.demand < 0 {
				ed.demand = 0
			}
			fmt.Printf("Updating electricity demand changed by %d and is at %d\n", increase, ed.demand)
			time.Sleep(time.Second)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fuelConsumption := &FuelConsumption{rateOfConsumption: 2}
	generatorSpeed := &GeneratorSpeed{maxSpeed: 5}
	electricityDemand := &ElectricityDemand{demand: 100}

	// Set up cross-references
	fuelConsumption.generatorSpeed = generatorSpeed
	generatorSpeed.fuelConsumption = fuelConsumption

	done := make(chan bool)

	go fuelConsumption.run(done)
	go generatorSpeed.run(done)
	go electricityDemand.run(done)

	fmt.Println("Starting simulation...")
	time.Sleep(20 * time.Second) // Run for 20 seconds

	close(done) // Signal all goroutines to stop
	time.Sleep(time.Second) // Give goroutines time to finish

	fmt.Println("Simulation complete.")
}