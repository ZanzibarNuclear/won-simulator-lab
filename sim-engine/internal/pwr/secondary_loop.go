package pwr

import (
	"fmt"

	"worldofnuclear.com/internal/simworks"
)

type SecondaryLoop struct {
	simworks.BaseComponent
	steamTemperature             float64 // in ˚C
	steamPressure                float64 // in MPa
	powerOperatedReliefValveOpen bool
	feedwaterPumpOn              bool
	feedwaterFlowRate            float64 // in m³/s
	feedheatersOn                bool
	feedwaterTemperatureOut      float64 // in Celsius; temperature of the feedwater as it enters the steam generator; related to efficiency of the steam generator
	feedwaterTemperatureIn       float64 // in ˚C; based on water leaving the condenser
}

func NewSecondaryLoop(name string, description string) *SecondaryLoop {
	return &SecondaryLoop{
		BaseComponent:           *simworks.NewBaseComponent(name, description),
		steamTemperature:        Config["common"]["room_temperature"],
		steamPressure:           Config["common"]["atmospheric_pressure"],
		feedwaterTemperatureOut: Config["secondary_loop"]["base_feedwater_temperature"],
		feedwaterTemperatureIn:  Config["secondary_loop"]["base_feedwater_temperature"],
	}
}

func (sl *SecondaryLoop) SteamTemperature() float64 {
	return sl.steamTemperature
}

func (sl *SecondaryLoop) SteamTemperatureUnit() string {
	return "˚C"
}

func (sl *SecondaryLoop) SteamPressure() float64 {
	return sl.steamPressure
}

func (sl *SecondaryLoop) SteamPressureUnit() string {
	return "MPa"
}

func (sl *SecondaryLoop) PowerOperatedReliefValveOpen() bool {
	return sl.powerOperatedReliefValveOpen
}

func (sl *SecondaryLoop) FeedwaterPumpOn() bool {
	return sl.feedwaterPumpOn
}

func (sl *SecondaryLoop) FeedwaterFlowRate() float64 {
	return sl.feedwaterFlowRate
}

func (sl *SecondaryLoop) FeedwaterFlowRateUnit() string {
	return "m³/s"
}

func (sl *SecondaryLoop) FeedheatersOn() bool {
	return sl.feedheatersOn
}

func (sl *SecondaryLoop) FeedwaterTemperatureOut() float64 {
	return sl.feedwaterTemperatureOut
}

func (sl *SecondaryLoop) FeedwaterTemperatureIn() float64 {
	return sl.feedwaterTemperatureIn
}

func (sl *SecondaryLoop) FeedwaterTemperatureUnit() string {
	return "˚C"
}

func (sl *SecondaryLoop) Status() map[string]interface{} {
	return map[string]interface{}{
		"about":                        sl.BaseComponent.Status(),
		"steamTemperature":             sl.SteamTemperature(),
		"steamTemperatureUnit":         sl.SteamTemperatureUnit(),
		"steamPressure":                sl.SteamPressure(),
		"steamPressureUnit":            sl.SteamPressureUnit(),
		"powerOperatedReliefValveOpen": sl.PowerOperatedReliefValveOpen(),
		"feedwaterPumpOn":              sl.FeedwaterPumpOn(),
		"feedwaterFlowRate":            sl.FeedwaterFlowRate(),
		"feedwaterFlowRateUnit":        sl.FeedwaterFlowRateUnit(),
		"feedheatersOn":                sl.FeedheatersOn(),
		"feedwaterTemperatureOut":      sl.FeedwaterTemperatureOut(),
		"feedwaterTemperatureIn":       sl.FeedwaterTemperatureIn(),
		"feedwaterTemperatureUnit":     sl.FeedwaterTemperatureUnit(),
	}
}

func (sl *SecondaryLoop) Print() {
	fmt.Printf("=> Secondary Loop\n")
	sl.BaseComponent.Print()
	fmt.Printf("Steam temperature: %.2f %s\n", sl.SteamTemperature(), sl.SteamTemperatureUnit())
	fmt.Printf("Steam pressure: %.2f %s\n", sl.SteamPressure(), sl.SteamPressureUnit())
	fmt.Printf("Power operated relief valve open: %v\n", sl.PowerOperatedReliefValveOpen())
	fmt.Printf("Feedwater pump on: %v\n", sl.FeedwaterPumpOn())
	fmt.Printf("Feedwater flow rate: %.2f %s\n", sl.FeedwaterFlowRate(), sl.FeedwaterFlowRateUnit())
	fmt.Printf("Feedheaters on: %v\n", sl.FeedheatersOn())
	fmt.Printf("Feedwater temperature out: %.2f %s\n", sl.FeedwaterTemperatureOut(), sl.FeedwaterTemperatureUnit())
	fmt.Printf("Feedwater temperature in: %.2f %s\n", sl.FeedwaterTemperatureIn(), sl.FeedwaterTemperatureUnit())
}

// Notes:
//
// Starting from the steam generator, water is heated, converted to steam,
// and the steam is heated more. Hot steam flows to the steam turbine,
// which spins, converting heat energy to mechanical energy. The turbine
// is connected to a generator, which produces electricity.
//
// Once the water leaves the turbine, the remaining heat is taken out as
// waste heat by condensers. A third cooling loop gets involved to
// remove the heat to outside sinks: a large body of water, cooling towers, etc.
//
// Water in this loop is returned to the steam generator by the feedwater pump.
// The water needs to be continually topped off to make up for evaporation
// and blowdown. Also, heating the water before it enters the  steam generator
// improves efficiency. Feedheaters take care of that.
//
// Most of the circulation is driven by natural convection. No pumping required;
// just lots of heat.
//
// Last but not least, safety valves prevent explosions due to excessive pressure.
// A little steam is vented when the pressure reach the safety threshold.

func (sl *SecondaryLoop) Update(s *simworks.Simulator) (map[string]interface{}, error) {
	// TODO: react to Steam Generator; determine steam temperature and pressure
	// steam moves at 60 mph during operation
	// if sl.steamTemperature < Config {
	// 	sl.steamTemperature += 10.0 // temperature increases some amount TODO: base this on Steam Generator
	// 	sl.steamPressure += 1.0     // pressure increases accordingly TODO: base this on steam temperature
	// }

	// vent steam when pressure is too high
	// if sl.steamPressure > Config["secondary_loop"]["mssv_pressure_threshold"] {
	// 	sl.mainSteamSafetyValveOpened = true
	// 	sl.steamPressure = Config["secondary_loop"]["mssv_pressure_threshold"] - 1.5 // TODO: research how much pressure would drop
	// 	sl.steamTemperature -= 30.0                                                  // temperature drops some amount TODO: research how much per event
	// } else {
	// 	sl.mainSteamSafetyValveOpened = false
	// }

	// if env.PowerOn {
	// 	// TODO: figure out less awkward way to adjust sub-components
	// 	if sl.openPowerOperatedReliefValve {
	// 		sl.steamPressure = sl.targetSteamPressure
	// 		sl.powerOperatedReliefValveOpened = true
	// 		sl.openPowerOperatedReliefValve = false
	// 	} else {
	// 		sl.powerOperatedReliefValveOpened = false
	// 	}

	// adjust feedwater temperature as needed
	// 	if sl.feedheatersOn && sl.feedwaterTemperature < Config["secondary_loop"]["heated_feedwater_temperature"] {
	// 		// increase temperature by 10 degree per minute until target is reached
	// 		sl.feedwaterTemperature += math.Min(Config["secondary_loop"]["heated_feedwater_temperature"]-sl.feedwaterTemperature, 0.5)
	// 	} else if !sl.feedheatersOn && sl.feedwaterTemperature > Config["secondary_loop"]["base_feedwater_temperature"] {
	// 		// decrease temperature by 10 degree per minute until base is reached
	// 		sl.feedwaterTemperature -= math.Min(sl.feedwaterTemperature-Config["secondary_loop"]["base_feedwater_temperature"], 0.5)
	// 	}
	// } else {
	// 	sl.SwitchOffFeedwaterPump()
	// 	sl.SwitchOffFeedheaters()
	// }

	return sl.Status(), nil
}

func (sl *SecondaryLoop) OpenPowerOperatedReliefValue(targetPressure float64) {
}

func (sl *SecondaryLoop) EmergencyMSSVReleased() bool {
	return true
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
