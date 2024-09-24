package pwr

import (
	"math"

	"worldofnuclear.com/internal/simworks"
)

type PrimaryLoop struct {
	simworks.BaseComponent
	pumpOn                   bool
	pumpPressure             float64 // in MPa
	flowRate                 float64 // in m³/s
	pumpHeat                 float64 // in MW
	boronConcentration       float64 // in parts per million (ppm)
	boronConcentrationTarget float64 // in parts per million (ppm)
	hotLegTemperature        float64 // in Celsius
	coldLegTemperature       float64 // in Celsius
}

func NewPrimaryLoop(name string, description string) *PrimaryLoop {
	return &PrimaryLoop{
		BaseComponent:            *simworks.NewBaseComponent(name, description),
		pumpOn:                   false,
		pumpPressure:             0.0,
		flowRate:                 0.0,
		pumpHeat:                 0.0,
		boronConcentration:       0.0,
		boronConcentrationTarget: 0.0,
		hotLegTemperature:        Config["common"]["room_temperature"],
		coldLegTemperature:       Config["common"]["room_temperature"],
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

func (pl *PrimaryLoop) BoronConcentrationTarget() float64 {
	return pl.boronConcentrationTarget
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
		"about":                    pl.BaseComponent.Status(),
		"pumpOn":                   pl.PumpOn(),
		"pumpPressure":             pl.PumpPressure(),
		"pumpPressureUnit":         pl.PumpPressureUnit(),
		"pumpHeat":                 pl.PumpHeat(),
		"pumpHeatUnit":             pl.PumpHeatUnit(),
		"flowRate":                 pl.FlowRate(),
		"flowRateUnit":             pl.FlowRateUnit(),
		"boronConcentration":       pl.BoronConcentration(),
		"boronConcentrationTarget": pl.BoronConcentrationTarget(),
		"boronConcentrationUnit":   pl.BoronConcentrationUnit(),
		"hotLegTemperature":        pl.HotLegTemperature(),
		"hotLegTemperatureUnit":    pl.TemperatureUnit(),
		"coldLegTemperature":       pl.ColdLegTemperature(),
		"coldLegTemperatureUnit":   pl.TemperatureUnit(),
	}
}

func (pl *PrimaryLoop) Update(s *simworks.Simulator) (map[string]interface{}, error) {
	pl.BaseComponent.Update(s)

	// TODO: use events to trigger changes - only need to set these values once per change
	if pl.pumpOn {
		// keep it simple for now. on full or off.
		pl.pumpPressure = Config["primary_loop"]["pump_on_pressure"]
		pl.pumpHeat = Config["primary_loop"]["pump_on_heat"]
		pl.flowRate = Config["primary_loop"]["pump_on_flow_rate"]

	} else {
		// no pressure, no flow, no change to boron concentration
		pl.pumpPressure = Config["primary_loop"]["pump_off_pressure"]
		pl.flowRate = Config["primary_loop"]["pump_off_flow_rate"]
		pl.pumpHeat = Config["primary_loop"]["pump_off_heat"]
	}

	// TODO: use events to trigger changes
	if pl.flowRate > 0 {
		// adjust boron concentration as needed
		if pl.boronConcentrationTarget != pl.boronConcentration {
			pl.boronConcentration = pl.boronConcentration + math.Copysign(
				math.Min(
					Config["primary_loop"]["max_boron_rate_of_change"],
					math.Abs(pl.boronConcentrationTarget-pl.boronConcentration),
				),
				pl.boronConcentrationTarget-pl.boronConcentration,
			)
		}
	}

	// Hot and cold leg temperatures depend on the core and steam generators.
	// Hot leg temperature depends on heat output of core reactor and primary loop pump,
	//   as well as how much water passes through the core.
	// Cold leg temperature depends on steam generator heat transfer

	// A good way to simplify might be to assume target values for each under
	// "normal" conditions.

	return pl.Status(), nil
}

// TODO: adopt Event-driven approach to system changes

// func (pl *PrimaryLoop) SwitchOnPump() {
// 	pl.pumpOn = true
// }

// func (pl *PrimaryLoop) SwitchOffPump() {
// 	pl.pumpOn = false
// }

// // set target amount in ppm
// // the system will approach this target concentration over time
// // no change happens while pump is off
// func (pl *PrimaryLoop) AdjustBoronConcentrationTarget(amount float64) {
// 	if amount < 0 {
// 		fmt.Printf("Boron concentration cannot be negative. You requested %f %s.\n", amount, pl.BoronConcentrationUnit())
// 		return
// 	}
// 	pl.boronConcentrationTarget = amount
// }
