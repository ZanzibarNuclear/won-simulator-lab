package sim

import (
	"fmt"
	"math"
)

type Pressurizer struct {
	BaseComponent
	pressure                float64
	temperature             float64
	heaterOn                bool
	heaterTargetTemperature float64
	targetPressure          float64
	heaterPower             float64 // in kW
	sprayNozzleOpen         bool
	reliefValveOpen         bool
	sprayFlowRate           float64 // in kg/s
	reliefValveFlow         float64 // in kg/s
}

func NewPressurizer(name string) *Pressurizer {
	return &Pressurizer{
		BaseComponent:           BaseComponent{Name: name},
		pressure:                15.5, // MPa, typical PWR pressurizer pressure
		temperature:             345,  // °C, typical PWR pressurizer temperature
		heaterOn:                false,
		heaterTargetTemperature: 345,  // °C, typical pressurizer temperature
		targetPressure:          15.5, // MPa, typical pressurizer pressure
		heaterPower:             1500, // kW, typical pressurizer heater capacity
		sprayNozzleOpen:         false,
		sprayFlowRate:           10, // kg/s, typical spray flow rate
		reliefValveOpen:         false,
		reliefValveFlow:         50, // kg/s, typical relief valve flow rate
	}
}

func (p *Pressurizer) GetName() string {
	return p.BaseComponent.Name
}

func (p *Pressurizer) Update(env *Environment, s *Simulation) {
	dt := 1.0 // always incrementing in 1-minute intervals (TODO: verify generated calculations below)

	// Simple model for pressure and temperature changes
	if p.heaterOn {
		p.temperature += (p.heaterPower / 1000) * dt // Simplified heating
		p.pressure += 0.01 * dt                      // Simplified pressure increase
	}

	if p.sprayNozzleOpen {
		p.temperature -= (p.sprayFlowRate / 100) * dt // Simplified cooling
		p.pressure -= 0.005 * dt                      // Simplified pressure decrease
	}

	if p.reliefValveOpen {
		p.pressure -= (p.reliefValveFlow / 100) * dt // Simplified rapid pressure decrease
	}

	// Ensure pressure and temperature stay within realistic bounds
	p.pressure = math.Max(10, math.Min(p.pressure, 17))         // MPa
	p.temperature = math.Max(300, math.Min(p.temperature, 370)) // °C
}

func (p *Pressurizer) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":                    p.Name,
		"pressure":                p.pressure,
		"temperature":             p.temperature,
		"heaterOn":                p.heaterOn,
		"heaterTargetTemperature": p.heaterTargetTemperature,
		"targetPressure":          p.targetPressure,
		"heaterPower":             p.heaterPower,
		"sprayNozzleOpen":         p.sprayNozzleOpen,
		"reliefValveOpen":         p.reliefValveOpen,
		"sprayFlowRate":           p.sprayFlowRate,
		"reliefValveFlow":         p.reliefValveFlow,
	}
}

func (p *Pressurizer) PrintStatus() {
	fmt.Printf("Pressurizer: %s\n", p.Name)
	fmt.Printf("\tPressure: %f\n", p.pressure)
	fmt.Printf("\tTemperature: %f\n", p.temperature)
	fmt.Printf("\tHeater On: %t\n", p.heaterOn)
	fmt.Printf("\tHeater Target Temperature: %f\n", p.heaterTargetTemperature)
	fmt.Printf("\tTarget Pressure: %f\n", p.targetPressure)
	fmt.Printf("\tHeater Power: %f\n", p.heaterPower)
}

func (p *Pressurizer) GetPressure() float64 {
	return p.pressure
}

func (p *Pressurizer) GetTemperature() float64 {
	return p.temperature
}

func (p *Pressurizer) SwitchHeater(on bool) {
	p.heaterOn = on
	if !p.heaterOn {
		p.heaterPower = 0
	} else {
		p.heaterPower = 1500
	}
}

func (p *Pressurizer) SetHeaterTargetTemperature(temp float64) {
	p.heaterTargetTemperature = temp
}

func (p *Pressurizer) OpenSprayNozzle(open bool) {
	p.sprayNozzleOpen = open
}

func (p *Pressurizer) CloseSprayNozzle() {
	p.sprayNozzleOpen = false
}

func (p *Pressurizer) OpenReliefValve(open bool) {
	p.reliefValveOpen = open
}
