package pwr

import (
	"worldofnuclear.com/internal/simworks"
)

// event codes
const (
	Event_pl_pumpSwitch         = "primary_loop.cooling_pump.switch"
	Event_pl_boronConcentration = "primary_loop.cvcs.boron_concentration_target"
	Event_pr_targetPressure     = "pressurizer.target_pressure.set"
	Event_pr_heaterPower        = "pressurizer.heater_power"
	Event_pr_sprayNozzle        = "pressurizer.spray_nozzle"
)

func NewPumpSwitchEvent(on bool) simworks.Event {
	return simworks.NewImmediateEventBool(Event_pl_pumpSwitch, on)
}

func NewBoronConcentrationEvent(concentration float64) simworks.Event {
	return simworks.NewAdjustmentEvent(Event_pl_boronConcentration, concentration)
}

func NewTargetPressureEvent(pressure float64) simworks.Event {
	return simworks.NewAdjustmentEvent(Event_pr_targetPressure, pressure)
}

func NewHeaterPowerEvent(on bool) simworks.Event {
	return simworks.NewImmediateEventBool(Event_pr_heaterPower, on)
}

func NewSprayNozzleEvent(open bool) simworks.Event {
	return simworks.NewImmediateEventBool(Event_pr_sprayNozzle, open)
}
