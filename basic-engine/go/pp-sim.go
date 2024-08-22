package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type FuelConsumption struct {
	rateOfIncrease    int
	rateOfConsumption int
	consumed          int
	generatorSpeed    *GeneratorSpeed
}

type GeneratorSpeed struct {
	speed            int
	maxSpeed         int
	fuelConsumption  *FuelConsumption
}

type ElectricityDemand struct {
	demand int
}

func (fc *FuelConsumption) run(wg *sync.WaitGroup, done chan bool) {
	defer wg.Done()
	for {
		select {
		case <-done:
			return
		default:
			if !fc.generatorSpeed.isRunningAtMax() {
				fc.rateOfConsumption += fc.rateOfIncrease
			}
			fc.consumed += fc.rateOfConsumption
			fmt.Printf("Consuming %d units of fuel. Total consumption at %d.\n", fc.rateOfConsumption, fc.consumed)
			time.Sleep(time.Second)
		}
	}
}

func (gs *GeneratorSpeed) run(wg *sync.WaitGroup, done chan bool) {
	defer wg.Done()
	for {
		select {
		case <-done:
			return
		default:
			gs.speed = int(math.Min(float64(gs.maxSpeed), float64(gs.fuelConsumption.rateOfConsumption*500)))
			fmt.Printf("Adjusting generator to run at %d RPMs\n", gs.speed)
			time.Sleep(time.Second)
		}
	}
}

func (gs *GeneratorSpeed) isRunningAtMax() bool {
	return gs.speed >= gs.maxSpeed
}

func (ed *ElectricityDemand) run(wg *sync.WaitGroup, done chan bool) {
	defer wg.Done()
	for {
		select {
		case <-done:
			return
		default:
			increase := rand.Intn(11) - 5
			ed.demand += increase
			if ed.demand < 0 {
				ed.demand = 0
			}
			fmt.Printf("Updating electricity demand changed by %d and is at %d kW\n", increase, ed.demand)
			time.Sleep(time.Second)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fuelConsumption := &FuelConsumption{
		rateOfIncrease:    1,
		rateOfConsumption: 1,
		consumed:          0,
	}

	generatorSpeed := &GeneratorSpeed{
		speed:    0,
		maxSpeed: 3600,
	}

	fuelConsumption.generatorSpeed = generatorSpeed
	generatorSpeed.fuelConsumption = fuelConsumption

	electricityDemand := &ElectricityDemand{
		demand: 100,
	}

	var wg sync.WaitGroup
	done := make(chan bool)

	wg.Add(3)
	go fuelConsumption.run(&wg, done)
	go generatorSpeed.run(&wg, done)
	go electricityDemand.run(&wg, done)

	fmt.Println("Starting simulation...")
	time.Sleep(20 * time.Second)
	close(done)

	wg.Wait()
	fmt.Println("Simulation complete.")
}