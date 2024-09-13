package sim

import (
	"fmt"
	"math"
)

const MSSV_PRESSURE_THRESHOLD = 8.0          // in MPa; main steam safety value
const TARGET_STEAM_TEMPERATURE = 285.0       // in Celsius
const TARGET_FEEDWATER_TEMPERATURE = 80.0    // in Celsius
const FEEDWATER_TEMPERATURE_INCREMENT = 10.0 // in Celsius
const BASE_FEEDWATER_TEMPERATURE = 40.0      // in Celsius

type SecondaryLoop struct {
	BaseComponent
	steamTemperature               float64 // in Celsius
	steamPressure                  float64 // in MPa
	mainSteamSafetyValveOpened     bool
	openPowerOperatedReliefValve   bool
	powerOperatedReliefValveOpened bool
	targetSteamPressure            float64 // in MPa
	feedwaterPumpOn                bool
	feedwaterFlowRate              float64 // in m³/s
	feedheatersOn                  bool
	feedwaterTemperature           float64 // in Celsius; temperature of the feedwater as it enters the steam generator; related to efficiency of the steam generator
}

func NewSecondaryLoop(name string) *SecondaryLoop {
	return &SecondaryLoop{
		BaseComponent:        BaseComponent{Name: name},
		steamTemperature:     100.0,
		steamPressure:        1.0,
		feedwaterFlowRate:    2.0, // 2 m³/s, 120 per minute
		feedwaterTemperature: BASE_FEEDWATER_TEMPERATURE,
	}
}

// Notes:
// Heat energy spins the turbine, which turns the generator to produce electricity.
// The remaining heat is taken out as waste heat by condensers, which involve
// the third loop to cooling towers and outside water sources.
//
// The feedwater pump is needed to top off water in the steam generator
// to make up for evaporation and blowdown.
// safety valve is needed to prevent explosions due to excessive pressure.
// Most of the circulation is driven by natural convection.

func (sl *SecondaryLoop) Update(env *Environment, s *Simulation) {
	// TODO: react to Steam Generator; determine steam temperature and pressure
	// steam moves at 60 mph during operation
	if sl.steamTemperature < TARGET_STEAM_TEMPERATURE {
		sl.steamTemperature += 10.0 // temperature increases some amount TODO: base this on Steam Generator
		sl.steamPressure += 1.0     // pressure increases accordingly TODO: base this on steam temperature
	}

	// vent steam when pressure is too high
	if sl.steamPressure > MSSV_PRESSURE_THRESHOLD {
		sl.mainSteamSafetyValveOpened = true
		sl.steamPressure -= MSSV_PRESSURE_THRESHOLD - 1.5 // TODO: research how much pressure would drop
		sl.steamTemperature -= 30.0                       // temperature drops some amount TODO: research how much per event
	} else {
		sl.mainSteamSafetyValveOpened = false
	}

	if env.PowerOn {
		// TODO: figure out less awkward way to adjust sub-components
		if sl.openPowerOperatedReliefValve {
			sl.steamPressure = sl.targetSteamPressure
			sl.powerOperatedReliefValveOpened = true
			sl.openPowerOperatedReliefValve = false
		} else {
			sl.powerOperatedReliefValveOpened = false
		}

		// adjust feedwater temperature as needed
		if sl.feedheatersOn && sl.feedwaterTemperature < TARGET_FEEDWATER_TEMPERATURE {
			// increase temperature by 10 degree per minute until target is reached
			sl.feedwaterTemperature += math.Min(TARGET_FEEDWATER_TEMPERATURE-sl.feedwaterTemperature, FEEDWATER_TEMPERATURE_INCREMENT)
		} else if !sl.feedheatersOn && sl.feedwaterTemperature > BASE_FEEDWATER_TEMPERATURE {
			// decrease temperature by 10 degree per minute until base is reached
			sl.feedwaterTemperature -= math.Min(sl.feedwaterTemperature-BASE_FEEDWATER_TEMPERATURE, FEEDWATER_TEMPERATURE_INCREMENT)
		}
	} else {
		sl.SwitchOffFeedwaterPump()
		sl.SwitchOffFeedheaters()
	}
}

func (sl *SecondaryLoop) FeedwaterVolume() float64 {
	if sl.feedwaterPumpOn {
		return sl.feedwaterFlowRate * 60
	}
	return 0.0
}

func (sl *SecondaryLoop) TargetFeedwaterTemperature() float64 {
	return TARGET_FEEDWATER_TEMPERATURE
}

func (sl *SecondaryLoop) OpenPowerOperatedReliefValue(targetPressure float64) {
	sl.openPowerOperatedReliefValve = true
	sl.targetSteamPressure = targetPressure
}

func (sl *SecondaryLoop) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":                       sl.Name,
		"steamTemperature":           sl.steamTemperature,
		"steamPressure":              sl.steamPressure,
		"mainSteamSafetyValveOpened": sl.mainSteamSafetyValveOpened,
		"feedwaterTemperature":       sl.feedwaterTemperature,
		"feedwaterPumpOn":            sl.feedwaterPumpOn,
		"feedwaterVolume":            sl.FeedwaterVolume(),
		"feedwaterHeatersOn":         sl.feedheatersOn,
	}
}

func (sl *SecondaryLoop) PrintStatus() {
	fmt.Printf("Secondary Loop: %s\n", sl.Name)
	fmt.Printf("\tSteam Temperature: %.2f °C\n", sl.steamTemperature)
	fmt.Printf("\tSteam Pressure: %.2f MPa\n", sl.steamPressure)
	fmt.Printf("\tMain Steam Safety Valve Released: %t\n", sl.mainSteamSafetyValveOpened)
	fmt.Printf("\tFeedwater Temperature: %.2f °C\n", sl.feedwaterTemperature)
	fmt.Printf("\tFeedwater Pump: %s\n", boolToString(sl.feedwaterPumpOn))
	fmt.Printf("\tFeedwater Volume: %.2f m³/min\n", sl.FeedwaterVolume())
	fmt.Printf("\tFeedwater Heaters: %s\n", boolToString(sl.feedheatersOn))
}

func boolToString(b bool) string {
	if b {
		return "On"
	}
	return "Off"
}

func (sl *SecondaryLoop) EmergencyMSSVReleased() bool {
	return sl.mainSteamSafetyValveOpened
}

func (sl *SecondaryLoop) SwitchOnFeedwaterPump() {
	sl.feedwaterPumpOn = true
	sl.feedwaterFlowRate = 2.0 // Reset to default flow rate when switched on
}

func (sl *SecondaryLoop) SwitchOffFeedwaterPump() {
	sl.feedwaterPumpOn = false
	sl.feedwaterFlowRate = 0.0 // No flow when pump is off
}

func (sl *SecondaryLoop) AdjustFeedwaterFlowRate(rate float64) {
	if !sl.feedwaterPumpOn {
		fmt.Println("Cannot adjust flow rate. Feedwater pump is off.")
		return
	}
	if rate < 0 {
		fmt.Println("Flow rate cannot be negative. Setting to 0.")
		sl.feedwaterFlowRate = 0
	} else {
		sl.feedwaterFlowRate = rate
	}
}

func (sl *SecondaryLoop) SwitchOnFeedheaters() {
	sl.feedheatersOn = true
	// Optionally, we could add some logic here to gradually increase the temperature
	// of the feedwater over time, simulating the heating process.
	// For example:
	// sl.returnWaterTemperature += 10 // Increase temperature by 10 degrees
	// This would depend on how often this method is called and how we want to model the heating process.
}

func (sl *SecondaryLoop) SwitchOffFeedheaters() {
	sl.feedheatersOn = false
	// Similarly, we could add logic here to gradually decrease the temperature
	// of the feedwater over time, simulating the cooling process when heaters are off.
}
