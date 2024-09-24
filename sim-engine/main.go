package main

import (
	"fmt"
	"time"

	"worldofnuclear.com/internal/simworks"
)

func main() {
	sim := simworks.NewSimulator("Simulator Works", "The inner workings to support simulations")
	fmt.Printf("Starting SimTime: %v\n", sim.Clock.FormatNow())
	sim.RunForABit(0, 0, 1, 0)
	fmt.Printf("After 1 minute SimTime: %v\n", sim.Clock.FormatNow())
	sim.RunForABit(0, 6, 0, 0)
	fmt.Printf("After 6 hours SimTime: %v\n", sim.Clock.FormatNow())
	fmt.Println("Running for a few days")
	sim.RunForABit(2, 0, 0, 0)
	fmt.Printf("After 2 days SimTime: %v\n", sim.Clock.FormatNow())

	fmt.Println("Running for a few days in own thread")
	go sim.RunForABit(3, 9, 12, 42)

	time.Sleep(5 * time.Second)
	fmt.Printf("After 3 days, 9 hours, 12 minute, 42 second SimTime: %v\n", sim.Clock.FormatNow())
	fmt.Println("All done.")

}
