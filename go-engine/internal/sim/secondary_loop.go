package sim

import (
	"fmt"
)

type SecondaryLoop struct {
	BaseComponent
	steamTemperature     float64 // in Celsius
	steamPressure        float64 // in MPa
	steamSafetyValveOpen bool
	feedwaterTemperature float64 // in Celsius
	feedwaterPumpOn      bool
	feedwaterFlowRate    float64 // in m³/s
	feedwaterHeatersOn   bool
}

func NewSecondaryLoop(name string) *SecondaryLoop {
	return &SecondaryLoop{
		BaseComponent:        BaseComponent{Name: name},
		steamTemperature:     100.0, // Initial value, can be adjusted as needed
		steamPressure:        1.0,   // Initial value, can be adjusted as needed
		feedwaterTemperature: 80.0,  // Initial value, can be adjusted as needed
		feedwaterFlowRate:    2.0,   // 2 m³/s, 120 per minute
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
	// TODO: react to Steam Generator
	// TODO: react to changes in feedwater flow rate
	// TODO: detect pressure relief events: pressure and temperature should drop,
	//   and make a note of the event
}

func (sl *SecondaryLoop) FeedwaterVolume() float64 {
	return sl.feedwaterFlowRate * 60
}

func (sl *SecondaryLoop) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":                 sl.Name,
		"steamTemperature":     sl.steamTemperature,
		"steamPressure":        sl.steamPressure,
		"steamSafetyValveOpen": sl.steamSafetyValveOpen,
		"feedwaterTemperature": sl.feedwaterTemperature,
		"feedwaterPumpOn":      sl.feedwaterPumpOn,
		"feedwaterVolume":      sl.FeedwaterVolume(),
		"feedwaterHeatersOn":   sl.feedwaterHeatersOn,
	}
}

func (sl *SecondaryLoop) PrintStatus() {
	fmt.Printf("Secondary Loop: %s\n", sl.Name)
	fmt.Printf("\tSteam Temperature: %.2f °C\n", sl.steamTemperature)
	fmt.Printf("\tSteam Pressure: %.2f MPa\n", sl.steamPressure)
	fmt.Printf("\tSteam Safety Valve: %s\n", boolToString(sl.steamSafetyValveOpen))
	fmt.Printf("\tFeedwater Temperature: %.2f °C\n", sl.feedwaterTemperature)
	fmt.Printf("\tFeedwater Pump: %s\n", boolToString(sl.feedwaterPumpOn))
	fmt.Printf("\tFeedwater Flow Rate: %.2f m³/min\n", sl.FeedwaterVolume())
	fmt.Printf("\tFeedwater Heaters: %s\n", boolToString(sl.feedwaterHeatersOn))
}

func boolToString(b bool) string {
	if b {
		return "On"
	}
	return "Off"
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

func (sl *SecondaryLoop) SwitchOnFeedwaterHeaters() {
	sl.feedwaterHeatersOn = true
	// Optionally, we could add some logic here to gradually increase the temperature
	// of the feedwater over time, simulating the heating process.
	// For example:
	// sl.returnWaterTemperature += 10 // Increase temperature by 10 degrees
	// This would depend on how often this method is called and how we want to model the heating process.
}

func (sl *SecondaryLoop) SwitchOffFeedwaterHeaters() {
	sl.feedwaterHeatersOn = false
	// Similarly, we could add logic here to gradually decrease the temperature
	// of the feedwater over time, simulating the cooling process when heaters are off.
}
