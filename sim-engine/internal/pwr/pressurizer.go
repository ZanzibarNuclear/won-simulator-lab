package pwr

import (
	"fmt"

	"worldofnuclear.com/internal/simworks"
)

type Pressurizer struct {
	simworks.BaseComponent
	targetPressure    float64
	pressure          float64
	temperature       float64
	heaterOn          bool
	heaterPower       float64 // in kW
	heaterTemperature float64
	sprayNozzleOpen   bool
	sprayFlowRate     float64 // in kg/s
}

func NewPressurizer(name string, description string) *Pressurizer {
	return &Pressurizer{
		BaseComponent:     *simworks.NewBaseComponent(name, description),
		targetPressure:    Config["pressurizer"]["target_pressure"],
		pressure:          Config["common"]["atmospheric_pressure"], // MPa, typical PWR pressurizer pressure
		temperature:       Config["common"]["room_temperature"],     // °C, typical PWR pressurizer temperature
		heaterTemperature: Config["common"]["room_temperature"],     // °C, typical pressurizer temperature
	}
}

// TODO: create event to change target pressure

func (p *Pressurizer) TargetPressure() float64 {
	return p.targetPressure
}

func (p *Pressurizer) Pressure() float64 {
	return p.pressure
}

func (p *Pressurizer) PressureUnit() string {
	return "MPa"
}

func (p *Pressurizer) Temperature() float64 {
	return p.temperature
}

func (p *Pressurizer) TemperatureUnit() string {
	return "°C"
}

func (p *Pressurizer) HeaterOn() bool {
	return p.heaterOn
}

func (p *Pressurizer) HeaterPower() float64 {
	return p.heaterPower
}

func (p *Pressurizer) HeaterTemperature() float64 {
	return p.heaterTemperature
}

func (p *Pressurizer) SprayNozzleOpen() bool {
	return p.sprayNozzleOpen
}

func (p *Pressurizer) SprayFlowRate() float64 {
	return p.sprayFlowRate
}

// TODO: turn this into an event
func (p *Pressurizer) SwitchOnHeater() {
	p.heaterOn = true
}

// TODO: turn this into an event
func (p *Pressurizer) SwitchOffHeater() {
	p.heaterOn = false
}

// TODO: turn this into an event
func (p *Pressurizer) OpenSprayNozzle() {
	p.sprayNozzleOpen = true
}

// TODO: turn this into an event
func (p *Pressurizer) CloseSprayNozzle() {
	p.sprayNozzleOpen = false
}

func (p *Pressurizer) Status() map[string]interface{} {
	return map[string]interface{}{
		"about":             p.BaseComponent.Status(),
		"targetPressure":    p.TargetPressure(),
		"pressure":          p.Pressure(),
		"pressureUnit":      p.PressureUnit(),
		"temperature":       p.Temperature(),
		"temperatureUnit":   p.TemperatureUnit(),
		"heaterOn":          p.HeaterOn(),
		"heaterPower":       p.HeaterPower(),
		"heaterTemperature": p.HeaterTemperature(),
		"sprayNozzleOpen":   p.SprayNozzleOpen(),
		"sprayFlowRate":     p.SprayFlowRate(),
	}
}

func (p *Pressurizer) Print() {
	fmt.Printf("=> Pressurizer\n")
	p.BaseComponent.Print()
	fmt.Printf("Target Pressure: %.2f %s\n", p.TargetPressure(), p.PressureUnit())
	fmt.Printf("Pressure: %.2f %s\n", p.Pressure(), p.PressureUnit())
	fmt.Printf("Temperature: %.2f %s\n", p.Temperature(), p.TemperatureUnit())
	fmt.Printf("Heater On: %t\n", p.HeaterOn())
	fmt.Printf("Heater Power: %.2f %s\n", p.HeaterPower(), p.TemperatureUnit())
	fmt.Printf("Heater Temperature: %.2f %s\n", p.HeaterTemperature(), p.TemperatureUnit())
	fmt.Printf("Spray Nozzle Open: %t\n", p.SprayNozzleOpen())
	fmt.Printf("Spray Flow Rate: %.2f %s\n", p.SprayFlowRate(), p.TemperatureUnit())
}

func (p *Pressurizer) Update(s *simworks.Simulator) {
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

	// TODO: turn this into an event
	// if p.pressure > RELIEF_VALVE_THRESHOLD_PRESSURE {
	// 	p.reliefValveOpened = true
	// 	p.pressure -= RELIEF_VALVE_DROP_RATE
	// } else {
	// 	p.reliefValveOpened = false
	// }

}
