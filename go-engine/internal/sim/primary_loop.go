package sim

import (
	"fmt"
	"math"
)

type PrimaryLoop struct {
	BaseComponent
	flowRate           float64 // in m³/s
	pumpOn             bool
	pumpPressure       float64 // in Pa
	boronConcentration float64 // in parts per million (ppm)
}

func NewPrimaryLoop(name string) *PrimaryLoop {
	return &PrimaryLoop{
		BaseComponent:      BaseComponent{Name: name},
		flowRate:           0,
		pumpOn:             false,
		pumpPressure:       0,
		boronConcentration: 0,
	}
}

func (pl *PrimaryLoop) Update(env *Environment, s *Simulation) {
	if pl.pumpOn {
		pl.pumpPressure = 1000000 // 1 MPa when pump is on
	} else {
		pl.pumpPressure = 0
	}

	// Simple flow rate calculation based on pump pressure
	// Assuming a linear relationship for simplicity
	pl.flowRate = math.Sqrt(pl.pumpPressure) * 0.001 // Arbitrary scaling factor
}

func (pl *PrimaryLoop) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":         pl.Name,
		"flowRate":     pl.flowRate,
		"pumpOn":       pl.pumpOn,
		"pumpPressure": pl.pumpPressure,
	}
}

func (pl *PrimaryLoop) PrintStatus() {
	fmt.Printf("Primary Loop: %s\n", pl.Name)
	fmt.Printf("\tFlow Rate: %.2f m³/s\n", pl.flowRate)
	fmt.Printf("\tPump Status: %v\n", pl.pumpOn)
	fmt.Printf("\tPump Pressure: %.2f Pa\n", pl.pumpPressure)
}

func (pl *PrimaryLoop) TogglePump() {
	pl.pumpOn = !pl.pumpOn
}

func (pl *PrimaryLoop) GetFlowRate() float64 {
	return pl.flowRate
}

func (pl *PrimaryLoop) AddBoron(amount float64) {
	pl.boronConcentration += amount
	if pl.boronConcentration < 0 {
		pl.boronConcentration = 0
	}
}

func (pl *PrimaryLoop) DiluteBoron(amount float64) {
	pl.boronConcentration -= amount
	if pl.boronConcentration < 0 {
		pl.boronConcentration = 0
	}
}
