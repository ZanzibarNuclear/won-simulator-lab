package pwr

import (
	"fmt"
)

// We are modeling the control banks that orchestrate across control rod assemblies.
// For reactivity, we need to control the position of the control rods. Although
// assemblies represent the physical groupings of control rods, we are modeling
// the control system that moves the rods within those assemblies.

// To move the core to criticality, first extract the control rods. Then move the
// control rods enough to see flux begin to rise exponentially, slightly super
// critical. At that point, put the control rods back a few positions, and that
// should do it.

// Not sure what to do with the gray banks yet. They should be useful for modeling
// "load following," whatever that is.

// Note: All banks of a given type move together at the moment. Need a targeted way
// to control the banks. Maybe later.

type ControlRods struct {
	controlBanks  [4]*ControlBank // full-strength absorbers; for power control during operation
	grayBanks     [2]*ControlBank // lower neutron absorption; for load following and fine reactivity control during operation
	shutdownBanks [4]*ControlBank // full-strength absorbers; for shutting down the core rapidly
}

func NewControlRods() *ControlRods {
	cr := &ControlRods{}

	// Initialize Control Banks
	cr.controlBanks[0] = NewControlBank("MA1", 4)
	cr.controlBanks[1] = NewControlBank("MA2", 4)
	cr.controlBanks[2] = NewControlBank("MB1", 4)
	cr.controlBanks[3] = NewControlBank("MB2", 4)

	cr.grayBanks[0] = NewControlBank("GR1", 8)
	cr.grayBanks[1] = NewControlBank("GR2", 8)

	cr.shutdownBanks[0] = NewControlBank("SD1", 8)
	cr.shutdownBanks[1] = NewControlBank("SD2", 8)
	cr.shutdownBanks[2] = NewControlBank("SD3", 8)
	cr.shutdownBanks[3] = NewControlBank("SD4", 8)

	return cr
}

func (cr *ControlRods) WithdrawShutdownBanks() {
	if cr.shutdownBanks[0].IsFullyWithdrawn() {
		return
	}

	for _, bank := range cr.shutdownBanks {
		bank.RaisePosition(int(Config["control_rods"]["withdrawal_rate"]))
	}
}

func (cr *ControlRods) ShutdownBanksFullyWithdrawn() bool {
	for _, bank := range cr.shutdownBanks {
		if !bank.IsFullyWithdrawn() {
			return false
		}
	}
	return true
}

func (cr *ControlRods) InsertShutdownBanks() {
	if cr.shutdownBanks[0].IsFullyInserted() {
		return
	}

	for _, bank := range cr.shutdownBanks {
		bank.LowerPosition(int(Config["control_rods"]["withdrawal_rate"]))
	}
}

func (cr *ControlRods) ShutdownBanksFullyInserted() bool {
	for _, bank := range cr.shutdownBanks {
		if !bank.IsFullyInserted() {
			return false
		}
	}
	return true
}

func (cr *ControlRods) AdjustControlBanks(target int) {
	if target < 0 || target > int(Config["control_rods"]["max_withdrawal_steps"]) {
		fmt.Printf("Tried to move control banks to %d, which is out of bounds\n", target)
		return
	}

	currentPosition := cr.controlBanks[0].Position()
	if currentPosition == target {
		return
	}
	for _, bank := range cr.controlBanks {
		if target > currentPosition {
			bank.RaisePosition(int(Config["control_rods"]["withdrawal_rate"]))
		} else {
			bank.LowerPosition(int(Config["control_rods"]["withdrawal_rate"]))
		}
	}
}

func (cr *ControlRods) AdjustGrayBanks(target int) {
	if target < 0 || target > int(Config["control_rods"]["max_withdrawal_steps"]) {
		fmt.Printf("Tried to move control banks to %d, which is out of bounds\n", target)
		return
	}

	currentPosition := cr.grayBanks[0].Position()
	if currentPosition == target {
		return
	}
	for _, bank := range cr.grayBanks {
		if target > currentPosition {
			bank.RaisePosition(int(Config["control_rods"]["withdrawal_rate"]))
		} else {
			bank.LowerPosition(int(Config["control_rods"]["withdrawal_rate"]))
		}
	}
}

// percent extraction
func (cr *ControlRods) AverageControlRodExtraction() float64 {
	// at the moment, all control rods will be in the same position
	return 1.0 - float64(cr.controlBanks[0].Position())/Config["control_rods"]["max_withdrawal_steps"]
}

func (cr *ControlRods) Scram() {
	for _, bank := range cr.controlBanks {
		bank.Scram()
	}
	for _, bank := range cr.grayBanks {
		bank.Scram()
	}
	for _, bank := range cr.shutdownBanks {
		bank.Scram()
	}
}

func (cr *ControlRods) Status() map[string]interface{} {
	status := make(map[string]interface{})

	for _, bank := range cr.controlBanks {
		status[bank.Label()] = bank.Status()
	}
	for _, bank := range cr.grayBanks {
		status[bank.Label()] = bank.Status()
	}
	for _, bank := range cr.shutdownBanks {
		status[bank.Label()] = bank.Status()
	}

	return status
}
