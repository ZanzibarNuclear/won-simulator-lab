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

func (sl *SecondaryLoop) SetSteamPressure(pressure float64) {
	sl.steamPressure = pressure
	sl.steamTemperature = InterpolateFromGivenPressure(pressure).Temperature
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
	sl.BaseComponent.Update(s)

	// TODO: try to move this to BaseComponent
	for i := range s.Events {
		event := s.Events[i]
		if event.IsPending() {
			if event.IsDue(s.CurrentMoment()) {
				event.SetInProgress()
			}
		}

		if event.IsInProgress() {
			if event.Immediate {
				sl.processInstantEvent(event)
			} else {
				sl.processGradualEvent(event)
			}
		}
	}

	// FIXME: react to Steam Generator; determine steam temperature and pressure
	// steam moves at 60 mph during operation

	// vent steam when pressure is too high
	if sl.steamPressure > Config["secondary_loop"]["mssv_pressure_threshold"] {
		s.QueueEvent(NewEvent_EmergencyMSSVReleased().ScheduleAt(s.CurrentMoment()))
		sl.steamPressure -= 0.5 // MPa; arbitrary value TODO: use a more realistic value
		newTemperature := InterpolateFromGivenPressure(sl.steamPressure)
		sl.steamTemperature = newTemperature.Temperature

		// Log the pressure release and temperature change
		fmt.Printf("Emergency MSSV released. Pressure dropped to %.2f MPa. Temperature adjusted to %.2f °C\n", sl.steamPressure, sl.steamTemperature)

		// FIXME: pressure has to drop; temperature, too
	}

	if sl.feedwaterPumpOn {
		sl.feedwaterFlowRate = CalcLinearIncrease(sl.feedwaterFlowRate, Config["secondary_loop"]["feedwater_flow_rate_target"], Config["secondary_loop"]["feedwater_flow_rate_step_up"])
	} else {
		sl.feedwaterFlowRate = CalcLinearDecrease(sl.feedwaterFlowRate, 0.0, Config["secondary_loop"]["feedwater_flow_rate_step_down"])
	}

	// automatic feedwater adjustments
	if sl.feedheatersOn {
		sl.feedwaterTemperatureOut = CalcLinearIncrease(sl.feedwaterTemperatureOut, Config["secondary_loop"]["heated_feedwater_temperature"], Config["secondary_loop"]["feedheater_step_up"])
	} else {
		sl.feedwaterTemperatureOut = CalcLinearDecrease(sl.feedwaterTemperatureOut, sl.feedwaterTemperatureIn, Config["secondary_loop"]["feedheater_step_down"])
	}

	if sl.powerOperatedReliefValveOpen {
		sl.steamPressure = CalcLinearDecrease(sl.steamPressure, Config["common"]["atmospheric_pressure"], 0.5)
	}

	sl.steamTemperature = InterpolateFromGivenPressure(sl.steamPressure).Temperature

	// TODO: deal with power outages here and in general

	return sl.Status(), nil
}

func (sl *SecondaryLoop) processInstantEvent(event *simworks.Event) {
	switch event.Code {
	case Event_sl_feedwaterPumpSwitch:
		sl.SwitchFeedwaterPump(event.Truthy())
	case Event_sl_feedheatersSwitch:
		sl.SwitchFeedheaters(event.Truthy())
	case Event_sl_powerOperatedReliefValve:
		sl.PowerOperatedReliefValve(event.Truthy())
	}
}

func (pl *SecondaryLoop) processGradualEvent(event *simworks.Event) {
	switch event.Code {
	default:
		return
	}
}

func (sl *SecondaryLoop) SwitchFeedwaterPump(on bool) {
	if on {
		sl.feedwaterPumpOn = true
	} else {
		sl.feedwaterPumpOn = false
	}
}

func (sl *SecondaryLoop) SwitchFeedheaters(on bool) {
	if on {
		sl.feedheatersOn = true
	} else {
		sl.feedheatersOn = false
	}
}

func (sl *SecondaryLoop) PowerOperatedReliefValve(open bool) {
	if open {
		sl.powerOperatedReliefValveOpen = true
	} else {
		sl.powerOperatedReliefValveOpen = false
	}
}
