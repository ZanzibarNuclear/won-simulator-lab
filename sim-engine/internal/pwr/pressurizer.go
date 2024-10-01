package pwr

import (
	"fmt"

	"worldofnuclear.com/internal/simworks"
)

type Pressurizer struct {
	simworks.BaseComponent
	pressure        float64
	temperature     float64
	heaterOn        bool
	heaterPower     float64 // in kW
	sprayNozzleOpen bool
	sprayFlowRate   float64 // in kg/s
}

func NewPressurizer(name string, description string) *Pressurizer {
	return &Pressurizer{
		BaseComponent: *simworks.NewBaseComponent(name, description),
		pressure:      Config["common"]["atmospheric_pressure"], // MPa, typical PWR pressurizer pressure
		temperature:   Config["common"]["room_temperature"],     // °C, typical PWR pressurizer temperature
	}
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

func (p *Pressurizer) HeaterOnHigh() bool {
	return p.heaterOn && p.heaterPower == Config["pressurizer"]["heater_high_power"]
}

func (p *Pressurizer) HeaterOnLow() bool {
	return p.heaterOn && p.heaterPower == Config["pressurizer"]["heater_low_power"]
}

func (p *Pressurizer) SprayNozzleOpen() bool {
	return p.sprayNozzleOpen
}

func (p *Pressurizer) SprayFlowRate() float64 {
	return p.sprayFlowRate
}

func (p *Pressurizer) Status() map[string]interface{} {
	return map[string]interface{}{
		"about":           p.BaseComponent.Status(),
		"pressure":        p.Pressure(),
		"pressureUnit":    p.PressureUnit(),
		"temperature":     p.Temperature(),
		"temperatureUnit": p.TemperatureUnit(),
		"heaterOn":        p.HeaterOn(),
		"heaterPower":     p.HeaterPower(),
		"sprayNozzleOpen": p.SprayNozzleOpen(),
		"sprayFlowRate":   p.SprayFlowRate(),
	}
}

func (p *Pressurizer) Print() {
	fmt.Printf("=> Pressurizer\n")
	p.BaseComponent.Print()
	fmt.Printf("Pressure: %.2f %s\n", p.Pressure(), p.PressureUnit())
	fmt.Printf("Temperature: %.2f %s\n", p.Temperature(), p.TemperatureUnit())
	fmt.Printf("Heater On: %t\n", p.HeaterOn())
	fmt.Printf("Heater Power: %.2f %s\n", p.HeaterPower(), p.TemperatureUnit())
	fmt.Printf("Spray Nozzle Open: %t\n", p.SprayNozzleOpen())
	fmt.Printf("Spray Flow Rate: %.2f %s\n", p.SprayFlowRate(), p.TemperatureUnit())
}

func (p *Pressurizer) Update(s *simworks.Simulator) (map[string]interface{}, error) {
	p.BaseComponent.Update(s)

	// TODO: try to move this to BaseComponent
	for _, event := range s.Events {
		if event.IsInProgress() {
			p.processEvent(event)
		}
	}

	if p.pressure > Config["pressurizer"]["relief_valve_threshold_pressure"] {
		reliefValveEvent := NewEvent_ReliefValveVent()
		reliefValveEvent.ScheduleAt(s.CurrentMoment())
		s.QueueEvent(reliefValveEvent)

		// Assume that venting lowers pressure and temperature along saturation curve
		newPressure := p.pressure - 1.0
		steamProperties := InterpolateFromGivenPressure(newPressure)
		fmt.Printf("Steam properties at %.2f MPa: Temperature %.2f °C, Pressure %.2f MPa\n",
			newPressure, steamProperties.Temperature, steamProperties.Pressure)
		p.pressure = steamProperties.Pressure
		p.temperature = steamProperties.Temperature
	}

	return p.Status(), nil
}

func (p *Pressurizer) processEvent(event *simworks.Event) {
	switch event.Code {
	case Event_pr_heaterPower:
		p.heaterOn = event.Truthy()
		if p.heaterOn {
			p.setHeaterLow()
		} else {
			p.heaterPower = 0.0
		}
		event.SetComplete()

	case Event_pr_sprayNozzle:
		p.sprayNozzleOpen = event.Truthy()
		if p.sprayNozzleOpen {
			p.sprayFlowRate = Config["pressurizer"]["spray_flow_rate"]
		} else {
			p.sprayFlowRate = 0.0
		}
		event.SetComplete()

	case Event_pr_targetPressure:
		targetValue := event.TargetValue
		p.adjustPressure(targetValue)
		if p.pressure >= targetValue {
			p.heaterPower = Config["pressurizer"]["heater_low_power"]
			event.SetComplete()
		}
	}
}

func (p *Pressurizer) setHeaterHigh() {
	p.heaterPower = Config["pressurizer"]["heater_high_power"]
}

func (p *Pressurizer) setHeaterLow() {
	p.heaterPower = Config["pressurizer"]["heater_low_power"]
}

func (p *Pressurizer) adjustPressure(targetValue float64) {

	// is heater on?
	if p.heaterOn {
		if p.pressure < targetValue {
			if !p.HeaterOnHigh() {
				p.setHeaterHigh()
			}
		} else {
			p.setHeaterLow()
		}

		if p.HeaterOnHigh() {
			p.temperature += 1.0
		}
	} else {
		p.temperature -= 1.0 // TODO: would be better to find out how quickly pressurizer cools when heaters are off
		if p.temperature < 290 {
			p.temperature = 290.0 // TODO: assuming cold leg temp
		}
	}

	if p.sprayNozzleOpen {
		p.temperature -= 5.0 // TODO: would be better to calculate using cold leg temp and volume of water sprayed
		if p.temperature < 290 {
			p.temperature = 290.0 // TODO: assuming cold leg temp
		}
	}

	// look up pressure from steam tables given temperature
	steamValues := InterpolateSteamProperties(p.temperature)
	p.pressure = steamValues.Pressure
}
