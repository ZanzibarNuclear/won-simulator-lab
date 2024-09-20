package sim

import (
	"fmt"
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
	sprayFlowRate     float64 // in kg/s
	reliefValveOpened bool
}

const TARGET_PRESSURE = 15.5                 // MPa, typical PWR pressurizer pressure
const TARGET_TEMPERATURE = 345.0             // °C, typical PWR pressurizer temperature
const HEATER_HIGH_POWER = 1500.0             // kW, typical pressurizer heater capacity
const HEATER_LOW_POWER = 50.0                // kW, enough to hold steady
const SPRAY_FLOW_RATE = 10.0                 // kg/s, typical spray flow rate
const RELIEF_VALVE_FLOW = 50.0               // kg/s, typical relief valve flow rate
const RELIEF_VALVE_THRESHOLD_PRESSURE = 17.0 // °C, typical PWR pressurizer temperature
const RELIEF_VALVE_DROP_DELTA = 2.5

func NewPressurizer(name string) *Pressurizer {
	return &Pressurizer{
		BaseComponent:     BaseComponent{Name: name},
		pressure:          0.0, // MPa, typical PWR pressurizer pressure
		targetPressure:    TARGET_PRESSURE,
		temperature:       ROOM_TEMPERATURE, // °C, typical PWR pressurizer temperature
		heaterPower:       0.0,              // kW, typical pressurizer heater capacity
		heaterTemperature: ROOM_TEMPERATURE, // °C, typical pressurizer temperature
		sprayFlowRate:     0.0,              // kg/s, typical spray flow rate
	}
}

func (p *Pressurizer) GetName() string {
	return p.BaseComponent.Name
}

func (p *Pressurizer) Update(env *Environment, s *Simulation) {
	// dt := 60 // seconds
	//
	// Simple formula for pressure changes:
	//   ΔP = (Q * β) / (V * Cp)
	//
	// Where:
	//   ΔP = Change in pressure
	//   Q = Heat input from the heater
	//   β = Coefficient of thermal expansion of water
	//   V = Volume of the pressurizer
	//   Cp = Specific heat capacity of water at constant pressure
	//
	// Also need the formula for water temperature based on pressure:
	//   T = T0 * (1 + β * ΔP)
	//
	// Where:
	//   T = Temperature of the water
	//   T0 = Initial temperature of the water
	//   β = Coefficient of thermal expansion of water
	//   ΔP = Change in pressure

	// TODO: this code is even simpler (and only directionally correct)
	if p.heaterOn {
		p.heaterTemperature = TARGET_TEMPERATURE
		if p.pressure < p.targetPressure || p.temperature < TARGET_TEMPERATURE {
			p.heaterPower = HEATER_HIGH_POWER

			// adjust pressure and temperature independently -- not realistic
			p.pressure += 1.0 // MPa, raise pressure by 0.5 MPa each update
			if p.pressure > p.targetPressure {
				p.pressure = p.targetPressure // cap pressure at TARGET_PRESSURE
			}
			p.temperature += 20.0
			if p.temperature > TARGET_TEMPERATURE {
				p.temperature = TARGET_TEMPERATURE
			}
		} else {
			p.heaterPower = HEATER_LOW_POWER // maintain pressure
		}
	} else {
		p.pressure -= 0.25 // MPa, assumption: pressure drops slowly when heater off
		if p.pressure < 0.0 {
			p.pressure = 0.0
		}
	}

	if p.sprayNozzleOpen {
		p.sprayFlowRate = SPRAY_FLOW_RATE
		p.pressure -= 0.5 // MPa; lower pressure
		p.temperature -= 20.0
		if p.temperature < ROOM_TEMPERATURE {
			p.temperature = ROOM_TEMPERATURE
		}
	} else {
		p.sprayFlowRate = 0.0
	}

	if p.pressure > RELIEF_VALVE_THRESHOLD_PRESSURE {
		p.reliefValveOpened = true
		p.pressure -= RELIEF_VALVE_DROP_DELTA
	} else {
		p.reliefValveOpened = false
	}

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
		"sprayFlowRate":     p.sprayFlowRate,
		"reliefValveOpened": p.reliefValveOpened,
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
	fmt.Printf("\tRelief Valve Opened: %t\n", p.reliefValveOpened)
}

func (p *Pressurizer) Pressure() float64 {
	return p.pressure
}

func (p *Pressurizer) SetTargetPressure(target float64) {
	p.targetPressure = target
}

func (p *Pressurizer) Temperature() float64 {
	return p.temperature
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
