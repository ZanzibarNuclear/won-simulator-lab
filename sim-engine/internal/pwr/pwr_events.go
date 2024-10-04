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
	Event_sl_emergencyMssvVent        = "secondary_loop.emergency_mssv.vented"
	Event_g_connectToGrid             = "generator.connect_to_grid"
	Event_g_connectToGridFailure      = "generator.connect_to_grid.failed"
	Event_rc_withdrawShutdownBanks    = "reactor_core.shutdown_banks.withdraw"
	Event_rc_insertShutdownBanks      = "reactor_core.shutdown_banks.insert"
	Event_rc_adjustControlRods        = "reactor_core.control_rods.adjust"
	Event_rc_adjustGrayRods           = "reactor_core.gray_rods.adjust"
	Event_rc_scram                    = "reactor_core.scram"
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

func NewEvent_PowerOperatedReliefValveSwitch(open bool) *simworks.Event {
	return simworks.NewImmediateEventBool(Event_sl_powerOperatedReliefValve, open)
}

func NewEvent_EmergencyMssvVent() *simworks.Event {
	return simworks.NewImmediateEvent(Event_sl_emergencyMssvVent)
}

func NewEvent_ConnectToGrid(connect bool) *simworks.Event {
	return simworks.NewImmediateEventBool(Event_g_connectToGrid, connect)
}

func NewEvent_ConnectToGridFailure() *simworks.Event {
	return simworks.NewImmediateEvent(Event_g_connectToGridFailure)
}

// TODO: expand event to handle adjustment event (over multiple steps) where the target is implied (e.g., fully out)
func NewEvent_WithdrawShutdownBanks() *simworks.Event {
	return simworks.NewAdjustmentEvent(Event_rc_withdrawShutdownBanks, 1.0)
}

func NewEvent_InsertShutdownBanks() *simworks.Event {
	return simworks.NewAdjustmentEvent(Event_rc_insertShutdownBanks, 1.0)
}

// TODO: expand event to hold map of details (e.g., which banks, which rods, how far out, etc.)
func NewEvent_AdjustControlRods(position float64) *simworks.Event {
	return simworks.NewAdjustmentEvent(Event_rc_adjustControlRods, position)
}

func NewEvent_AdjustGrayRods(position float64) *simworks.Event {
	return simworks.NewAdjustmentEvent(Event_rc_adjustGrayRods, position)
}

func NewEvent_Scram() *simworks.Event {
	return simworks.NewImmediateEvent(Event_rc_scram)
}
