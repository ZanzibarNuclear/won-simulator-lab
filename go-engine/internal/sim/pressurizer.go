package sim

import (
	"fmt"
	"math"
)

type Pressurizer struct {
	BaseComponent
	targetPressure    float64
	pressure          float64
	temperature       float64
	heaterOn          bool
	heaterPower       float64 // in kW
	heaterTemperature float64
	sprayNozzleOpen   bool
	reliefValveOpen   bool
	sprayFlowRate     float64 // in kg/s
	reliefValveFlow   float64 // in kg/s
}

const DEFAULT_TARGET_PRESSURE = 15.5         // MPa, typical PWR pressurizer pressure
const TARGET_TEMPERATURE = 345.0             // °C, typical PWR pressurizer temperature
const HEATER_HIGH_POWER = 1500.0             // kW, typical pressurizer heater capacity
const HEATER_LOW_POWER = 50.0                // kW, enough to hold steady
const SPRAY_FLOW_RATE = 10.0                 // kg/s, typical spray flow rate
const RELIEF_VALVE_FLOW = 50.0               // kg/s, typical relief valve flow rate
const RELIEF_VALUE_THRESHOLD_PRESSURE = 17.0 // °C, typical PWR pressurizer temperature

func NewPressurizer(name string) *Pressurizer {
	return &Pressurizer{
		BaseComponent:     BaseComponent{Name: name},
		targetPressure:    DEFAULT_TARGET_PRESSURE, // MPa, has default but adjustable
		pressure:          0.0,                     // MPa, typical PWR pressurizer pressure
		temperature:       ROOM_TEMPERATURE,        // °C, typical PWR pressurizer temperature
		heaterPower:       0.0,                     // kW, typical pressurizer heater capacity
		heaterTemperature: ROOM_TEMPERATURE,        // °C, typical pressurizer temperature
		sprayFlowRate:     0.0,                     // kg/s, typical spray flow rate
		reliefValveFlow:   0.0,                     // kg/s, typical relief valve flow rate
	}
}

func (p *Pressurizer) GetName() string {
	return p.BaseComponent.Name
}

// FIXME: redo when brain is calm
func (p *Pressurizer) Update(env *Environment, s *Simulation) {
	// dt := s.GetTimeStep()

	// Simple formula for pressure changes
	// ΔP = (Q * β) / (V * Cp)
	//
	//Where:
	// ΔP = Change in pressure
	// Q = Heat input from the heater
	// β = Coefficient of thermal expansion of water
	// V = Volume of the pressurizer
	// Cp = Specific heat capacity of water at constant pressure
	//
	// Also need the formula for water temperature based on pressure.
	// T = T0 * (1 + β * ΔP)
	//
	// Where:
	// T = Temperature of the water
	// T0 = Initial temperature of the water
	// β = Coefficient of thermal expansion of water
	// ΔP = Change in pressure

	
}

func (p *Pressurizer) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":              p.Name,
		"pressure":          p.pressure,
		"temperature":       p.temperature,
		"heaterOn":          p.heaterOn,
		"heaterTemperature": p.heaterTemperature,
		"targetPressure":    p.targetPressure,
		"heaterPower":       p.heaterPower,
		"sprayNozzleOpen":   p.sprayNozzleOpen,
		"reliefValveOpen":   p.reliefValveOpen,
		"sprayFlowRate":     p.sprayFlowRate,
		"reliefValveFlow":   p.reliefValveFlow,
	}
}

func (p *Pressurizer) PrintStatus() {
	fmt.Printf("Pressurizer: %s\n", p.Name)
	fmt.Printf("\tPressure: %f\n", p.pressure)
	fmt.Printf("\tTemperature: %f\n", p.temperature)
	fmt.Printf("\tHeater On: %t\n", p.heaterOn)
	fmt.Printf("\tHeater Temperature: %f\n", p.heaterTemperature)
	fmt.Printf("\tTarget Pressure: %f\n", p.targetPressure)
	fmt.Printf("\tHeater Power: %f\n", p.heaterPower)
	fmt.Printf("\tSpray Nozzle Open: %t\n", p.sprayNozzleOpen)
	fmt.Printf("\tSpray Flow Rate: %f\n", p.sprayFlowRate)
	fmt.Printf("\tRelief Valve Opened: %t\n", p.reliefValveOpen)
	fmt.Printf("\tRelief Valve Flow Rate: %f\n", p.reliefValveFlow)
}

func (p *Pressurizer) Pressure() float64 {
	return p.pressure
}

func (p *Pressurizer) Temperature() float64 {
	return p.temperature
}

func (p *Pressurizer) SetTargetPressure(target float64) {
	p.targetPressure = target
}

func (p *Pressurizer) SwitchOnHeater() {
	p.heaterOn = true
}

func (p *Pressurizer) SwitchOffHeater() {
	p.heaterOn = false
}

func (p *Pressurizer) OpenSprayNozzle() {
	p.sprayNozzleOpen = true
}

func (p *Pressurizer) CloseSprayNozzle() {
	p.sprayNozzleOpen = false
}
