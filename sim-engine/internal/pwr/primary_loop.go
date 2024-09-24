package pwr

import (
	"fmt"
	"math"

	"worldofnuclear.com/internal/simworks"
)

type PrimaryLoop struct {
	simworks.BaseComponent
	pumpOn             bool
	pumpPressure       float64 // in MPa
	flowRate           float64 // in m³/s
	pumpHeat           float64 // in MW
	boronConcentration float64 // in parts per million (ppm)
	hotLegTemperature  float64 // in Celsius
	coldLegTemperature float64 // in Celsius
}

func NewPrimaryLoop(name string, description string) *PrimaryLoop {
	return &PrimaryLoop{
		BaseComponent:      *simworks.NewBaseComponent(name, description),
		pumpOn:             false,
		pumpPressure:       0.0,
		flowRate:           0.0,
		pumpHeat:           0.0,
		boronConcentration: 0.0,
		hotLegTemperature:  Config["common"]["room_temperature"],
		coldLegTemperature: Config["common"]["room_temperature"],
	}
}

func (pl *PrimaryLoop) PumpOn() bool {
	return pl.pumpOn
}

func (pl *PrimaryLoop) PumpPressure() float64 {
	return pl.pumpPressure
}

func (pl *PrimaryLoop) PumpPressureUnit() string {
	return "MPa"
}

func (pl *PrimaryLoop) PumpHeat() float64 {
	return pl.pumpHeat
}

func (pl *PrimaryLoop) PumpHeatUnit() string {
	return "MW"
}

// Calculates flow volume per minute when pumps are on
func (pl *PrimaryLoop) FlowRate() float64 {
	return pl.flowRate
}

func (pl *PrimaryLoop) FlowRateUnit() string {
	return "m³/s"
}

func (pl *PrimaryLoop) BoronConcentration() float64 {
	return pl.boronConcentration
}

func (pl *PrimaryLoop) BoronConcentrationUnit() string {
	return "ppm"
}

func (pl *PrimaryLoop) HotLegTemperature() float64 {
	return pl.hotLegTemperature
}

func (pl *PrimaryLoop) TemperatureUnit() string {
	return "C"
}

func (pl *PrimaryLoop) ColdLegTemperature() float64 {
	return pl.coldLegTemperature
}

func (pl *PrimaryLoop) Status() map[string]interface{} {
	return map[string]interface{}{
		"about":                  pl.BaseComponent.Status(),
		"pumpOn":                 pl.PumpOn(),
		"pumpPressure":           pl.PumpPressure(),
		"pumpPressureUnit":       pl.PumpPressureUnit(),
		"pumpHeat":               pl.PumpHeat(),
		"pumpHeatUnit":           pl.PumpHeatUnit(),
		"flowRate":               pl.FlowRate(),
		"flowRateUnit":           pl.FlowRateUnit(),
		"boronConcentration":     pl.BoronConcentration(),
		"boronConcentrationUnit": pl.BoronConcentrationUnit(),
		"hotLegTemperature":      pl.HotLegTemperature(),
		"hotLegTemperatureUnit":  pl.TemperatureUnit(),
		"coldLegTemperature":     pl.ColdLegTemperature(),
		"coldLegTemperatureUnit": pl.TemperatureUnit(),
	}
}

func (pl *PrimaryLoop) Print() {
	fmt.Printf("=> Primary Loop\n")
	pl.BaseComponent.Print()
	fmt.Printf("Pump on: %v\n", pl.PumpOn())
	fmt.Printf("Pump pressure: %.2f %s\n", pl.PumpPressure(), pl.PumpPressureUnit())
	fmt.Printf("Pump heat: %.2f %s\n", pl.PumpHeat(), pl.PumpHeatUnit())
	fmt.Printf("Flow rate: %.2f %s\n", pl.FlowRate(), pl.FlowRateUnit())
	fmt.Printf("Boron concentration: %.2f %s\n", pl.BoronConcentration(), pl.BoronConcentrationUnit())
	fmt.Printf("Hot leg temperature: %.1f %s\n", pl.HotLegTemperature(), pl.TemperatureUnit())
	fmt.Printf("Cold leg temperature: %.1f %s\n", pl.ColdLegTemperature(), pl.TemperatureUnit())
}

func (pl *PrimaryLoop) Update(s *simworks.Simulator) (map[string]interface{}, error) {
	pl.BaseComponent.Update(s)

	// TODO: try to move this to BaseComponent
	for i := range s.Events {
		event := &s.Events[i]
		if event.IsPending() {
			if event.IsDue(s.CurrentMoment()) {
				event.SetInProgress()
			}
		}

		if event.IsInProgress() {
			if event.Immediate {
				pl.processInstantEvent(event)
			} else {
				pl.processGradualEvent(event)
			}
		}
	}

	// TODO: add updates based on environment and other components

	return pl.Status(), nil
}

func (pl *PrimaryLoop) processInstantEvent(event *simworks.Event) {
	switch event.Code {
	case Event_pl_pumpSwitch:
		pl.switchPump(event.Truthy())
		if pl.pumpOn {
			pl.pumpPressure = Config["primary_loop"]["pump_on_pressure"]
			pl.pumpHeat = Config["primary_loop"]["pump_on_heat"]
			pl.flowRate = Config["primary_loop"]["pump_on_flow_rate"]
		} else {
			pl.pumpPressure = Config["primary_loop"]["pump_off_pressure"]
			pl.pumpHeat = Config["primary_loop"]["pump_off_heat"]
			pl.flowRate = Config["primary_loop"]["pump_off_flow_rate"]
		}
	}
	event.SetComplete()
}

func (pl *PrimaryLoop) processGradualEvent(event *simworks.Event) {
	targetValue := event.TargetValue
	switch event.Code {
	case Event_pl_boronConcentration:
		pl.adjustBoron(targetValue)
		if pl.boronConcentration == targetValue {
			event.SetComplete()
			fmt.Printf("Boron concentration adjusted to %v\n", pl.boronConcentration)
		}
	}
}

func (pl *PrimaryLoop) adjustBoron(targetValue float64) {
	if pl.pumpOn {
		if pl.boronConcentration != targetValue {
			pl.boronConcentration = pl.boronConcentration + math.Copysign(
				math.Min(
					Config["primary_loop"]["max_boron_rate_of_change"],
					math.Abs(targetValue-pl.boronConcentration),
				),
				targetValue-pl.boronConcentration,
			)
		}
	}
}

func (pl *PrimaryLoop) switchPump(on bool) {
	pl.pumpOn = on
}
