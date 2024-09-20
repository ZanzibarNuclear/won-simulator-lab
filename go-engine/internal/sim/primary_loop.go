package sim

import (
	"fmt"
	"math"
)

type PrimaryLoop struct {
	BaseComponent
	pumpOn                   bool
	pumpPressure             float64 // in MPa
	flowRate                 float64 // in m³/s
	boronConcentration       float64 // in parts per million (ppm)
	boronConcentrationTarget float64 // in parts per million (ppm)
	pressurizer              *Pressurizer
}

// TODO: consider where to track hot / cold leg temperatures;
// hot leg is influenced by reactor core heat output;
// cold leg influenced by condenser output

// useful constants
// TODO: make some of these configurable
const PUMP_ON_PRESSURE = 1.0   // MPa
const PUMP_ON_FLOW_RATE = 20.0 // in m³/s
const PUMP_OFF_PRESSURE = 0
const PUMP_OFF_FLOW_RATE = 0
const MAX_BORON_RATE_OF_CHANGE = 5.0   // ppm/minute
const MAX_BORON_CONCENTRATION = 2500.0 // ppm

func NewPrimaryLoop(name string) *PrimaryLoop {
	return &PrimaryLoop{
		BaseComponent:            BaseComponent{Name: name},
		flowRate:                 0,
		pumpOn:                   false,
		pumpPressure:             0,
		boronConcentration:       0,
		boronConcentrationTarget: 0,
	}
}

func (pl *PrimaryLoop) Update(env *Environment, s *Simulation) {

	if pl.pumpOn {
		// TODO: does it make sense to allow variable pump speed?
		// then, perhaps half pressure leads to half the flow volume

		// keep it simple for now. on full or off.
		pl.pumpPressure = PUMP_ON_PRESSURE
		pl.flowRate = PUMP_ON_FLOW_RATE

		// adjust boron concentration as needed
		if pl.boronConcentrationTarget != pl.boronConcentration {
			pl.boronConcentration = pl.boronConcentration + math.Copysign(
				math.Min(
					MAX_BORON_RATE_OF_CHANGE,
					math.Abs(pl.boronConcentrationTarget-pl.boronConcentration),
				),
				pl.boronConcentrationTarget-pl.boronConcentration,
			)
		}
	} else {
		// no pressure, no flow, no change to boron concentration
		pl.pumpPressure = PUMP_OFF_PRESSURE
		pl.flowRate = PUMP_OFF_FLOW_RATE
	}

}

// Returns the current pump pressure in Pa
func (pl *PrimaryLoop) Pressure() float64 {
	return pl.pumpPressure
}

func (pl *PrimaryLoop) PressureUnit() string {
	return "MPa"
}

func (pl *PrimaryLoop) FlowVolumeUnit() string {
	return "m³/min"
}

// Calculates flow volume per minute when pumps are on
func (pl *PrimaryLoop) FlowVolume() float64 {
	if pl.pumpOn {
		return pl.flowRate * 60.0
	} else {
		// assumes water stops after 1 minute of pumps switching off
		return 0.0
	}
}

func (pl *PrimaryLoop) BoronConcentration() float64 {
	return pl.boronConcentration
}

func (pl *PrimaryLoop) BoronConcentrationUnit() string {
	return "ppm"
}

func (pl *PrimaryLoop) BoronConcentrationTarget() float64 {
	return pl.boronConcentrationTarget
}

func (pl *PrimaryLoop) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":                     pl.Name,
		"pumpOn":                   pl.pumpOn,
		"pumpPressure":             pl.Pressure(),
		"pressureUnit":             pl.PressureUnit(),
		"flowVolume":               pl.FlowVolume(),
		"flowVolumeUnit":           pl.FlowVolumeUnit(),
		"boronConcentration":       pl.BoronConcentration(),
		"boronConcentrationTarget": pl.BoronConcentrationTarget(),
		"boronConcentrationUnit":   pl.BoronConcentrationUnit(),
	}
}

func (pl *PrimaryLoop) PrintStatus() {
	fmt.Printf("Primary Loop: %s\n", pl.Name)
	pumpStatus := "Off"
	if pl.pumpOn {
		pumpStatus = "On"
	}
	fmt.Printf("\tPump Status: %s\n", pumpStatus)
	fmt.Printf("\tPump Pressure: %.2f %s\n", pl.pumpPressure, pl.PressureUnit())
	fmt.Printf("\tFlow Volume: %.2f %s\n", pl.FlowVolume(), pl.FlowVolumeUnit())
	fmt.Printf("\tBoron Concentration: %.2f %s\n", pl.BoronConcentration(), pl.BoronConcentrationUnit())
	fmt.Printf("\tBoron Concentration Target: %.2f %s\n", pl.BoronConcentrationTarget(), pl.BoronConcentrationUnit())
}

func (pl *PrimaryLoop) SwitchOnPump() {
	pl.pumpOn = true
}

func (pl *PrimaryLoop) SwitchOffPump() {
	pl.pumpOn = false
}

// set target amount in ppm
// the system will approach this target concentration over time
// no change happens while pump is off
func (pl *PrimaryLoop) AdjustBoronConcentrationTarget(amount float64) {
	if amount < 0 {
		fmt.Printf("Boron concentration cannot be negative. You requested %f %s.\n", amount, pl.BoronConcentrationUnit())
		return
	}
	pl.boronConcentrationTarget = amount
}
