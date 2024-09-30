package pwr

import (
	"worldofnuclear.com/internal/simworks"
)

// event codes
const (
	Event_pl_pumpSwitch               = "primary_loop.cooling_pump.switch"
	Event_pl_boronConcentration       = "primary_loop.cvcs.boron_concentration_target"
	Event_pr_targetPressure           = "pressurizer.target_pressure.set"
	Event_pr_heaterPower              = "pressurizer.heater_power"
	Event_pr_sprayNozzle              = "pressurizer.spray_nozzle"
	Event_pr_reliefValveVent          = "pressurizer.relief_valve.vented"
	Event_sl_feedwaterPumpSwitch      = "secondary_loop.feedwater_pump.switch"
	Event_sl_feedheatersSwitch        = "secondary_loop.feedheaters.switch"
	Event_sl_powerOperatedReliefValve = "secondary_loop.power_operated_relief_valve.switch"
	Event_sl_emergencyMSSVReleased    = "secondary_loop.emergency_mssv.released"
)

func NewEvent_PumpSwitch(on bool) *simworks.Event {
	return simworks.NewImmediateEventBool(Event_pl_pumpSwitch, on)
}

func NewEvent_BoronConcentration(concentration float64) *simworks.Event {
	return simworks.NewAdjustmentEvent(Event_pl_boronConcentration, concentration)
}

func NewEvent_HeaterPower(on bool) *simworks.Event {
	return simworks.NewImmediateEventBool(Event_pr_heaterPower, on)
}

func NewEvent_SprayNozzle(open bool) *simworks.Event {
	return simworks.NewImmediateEventBool(Event_pr_sprayNozzle, open)
}

func NewEvent_TargetPressure(pressure float64) *simworks.Event {
	return simworks.NewAdjustmentEvent(Event_pr_targetPressure, pressure)
}

func NewEvent_ReliefValveVent() *simworks.Event {
	return simworks.NewImmediateEvent(Event_pr_reliefValveVent)
}

func NewEvent_FeedwaterPumpSwitch(on bool) *simworks.Event {
	return simworks.NewImmediateEventBool(Event_sl_feedwaterPumpSwitch, on)
}

func NewEvent_FeedheatersSwitch(on bool) *simworks.Event {
	return simworks.NewImmediateEventBool(Event_sl_feedheatersSwitch, on)
}

func NewEvent_PowerOperatedReliefValveSwitch(on bool) *simworks.Event {
	return simworks.NewImmediateEventBool(Event_sl_powerOperatedReliefValve, on)
}

func NewEvent_EmergencyMSSVReleased() *simworks.Event {
	return simworks.NewImmediateEvent(Event_sl_emergencyMSSVReleased)
}
