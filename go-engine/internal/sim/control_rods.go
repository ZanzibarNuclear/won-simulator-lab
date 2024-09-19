package sim

import (
	"fmt"
)

const (
	MAX_WITHDRAWAL_STEPS = 200
	CONTROL_ROD_SPEED    = 40 // steps per minute
)

type ControlBank struct {
	position         int // in steps
	numRods          int
	withdrawalTarget int // position to move toward
}

func (cb *ControlBank) SetWithdrawalTarget(target int) {
	fmt.Printf("Setting withdrawal target for bank to %d\n", target)
	if target < 0 {
		cb.withdrawalTarget = 0
	} else if target > MAX_WITHDRAWAL_STEPS {
		cb.withdrawalTarget = MAX_WITHDRAWAL_STEPS
	} else {
		cb.withdrawalTarget = target
	}
}

func (cb *ControlBank) Withdrawn() bool {
	return cb.position == MAX_WITHDRAWAL_STEPS
}

type ControlRods struct {
	controlBanks          [7]ControlBank
	shutdownBanks         [4]ControlBank
	withdrawShutdownBanks bool // gradually
	insertShutdownBanks   bool // gradually
	scramNow              bool // drop all rods at once
}

func NewControlRods() *ControlRods {
	cr := &ControlRods{}

	for i := 0; i < 5; i++ {
		cr.controlBanks[i] = ControlBank{numRods: 4}
	}
	cr.controlBanks[5] = ControlBank{numRods: 8}
	cr.controlBanks[6] = ControlBank{numRods: 9}

	for i := 0; i < 4; i++ {
		cr.shutdownBanks[i] = ControlBank{numRods: 8}
	}

	return cr
}

func (cr *ControlRods) Update() {
	if cr.withdrawShutdownBanks && !cr.ShutdownBanksFullyWithdrawn() {
	Test1:
		for i, bank := range cr.shutdownBanks {
			if !bank.Withdrawn() {
				fmt.Printf("Setting withdrawal target for bank %d\n", i)
				cr.shutdownBanks[i].SetWithdrawalTarget(MAX_WITHDRAWAL_STEPS)
				break Test1 // Only start withdrawing one bank at a time
			}
		}
	}

	if cr.insertShutdownBanks && !cr.ShutdownBanksFullyInserted() {
		fmt.Printf("Inserting shutdown banks\n")
		// Insert in reverse order from withdrawal
	Test2:
		for i := len(cr.shutdownBanks) - 1; i >= 0; i-- {
			bank := cr.shutdownBanks[i]
			if bank.position > 0 {
				fmt.Printf("Setting withdrawal target for bank %d\n", i)
				bank.SetWithdrawalTarget(0)
				break Test2
			}
		}
	}

	for i, bank := range cr.controlBanks {
		if bank.withdrawalTarget > bank.position {
			cr.controlBanks[i].position += min(CONTROL_ROD_SPEED, bank.withdrawalTarget-bank.position) // Use index to update
		} else if bank.withdrawalTarget < bank.position {
			cr.controlBanks[i].position -= min(CONTROL_ROD_SPEED, bank.position-bank.withdrawalTarget)
		}
	}

	for i, bank := range cr.shutdownBanks {
		if bank.withdrawalTarget > bank.position {
			cr.shutdownBanks[i].position += min(CONTROL_ROD_SPEED, bank.withdrawalTarget-bank.position)
		} else if bank.withdrawalTarget < bank.position {
			cr.shutdownBanks[i].position -= min(CONTROL_ROD_SPEED, bank.position-bank.withdrawalTarget)
		}
	}
}

func (cr *ControlRods) Scram() {
	cr.scramNow = true
}

func (cr *ControlRods) InitiateShutdownBankWithdrawal() {
	cr.withdrawShutdownBanks = true
	cr.insertShutdownBanks = false
}

func (cr *ControlRods) InitiateShutdownBankInsertion() {
	cr.insertShutdownBanks = true
	cr.withdrawShutdownBanks = false
}

func (cr *ControlRods) ShutdownBanksFullyWithdrawn() bool {
	for _, bank := range cr.shutdownBanks {
		if !bank.Withdrawn() {
			return false
		}
	}
	return true
}

func (cr *ControlRods) ShutdownBanksFullyInserted() bool {
	for _, bank := range cr.shutdownBanks {
		if bank.position > 0 {
			return false
		}
	}
	return true
}

func (cr *ControlRods) Status() map[string]interface{} {
	controlBanksStatus := make([]map[string]interface{}, len(cr.controlBanks))
	for i, bank := range cr.controlBanks {
		controlBanksStatus[i] = map[string]interface{}{
			"controlBankNum":   i,
			"position":         bank.position,
			"withdrawalTarget": bank.withdrawalTarget,
		}
	}

	shutdownBanksStatus := make([]map[string]interface{}, len(cr.shutdownBanks))
	for i, bank := range cr.shutdownBanks {
		shutdownBanksStatus[i] = map[string]interface{}{
			"shutdownBankNum":  i,
			"position":         bank.position,
			"withdrawalTarget": bank.withdrawalTarget,
		}
	}

	return map[string]interface{}{
		"controlBanks":          controlBanksStatus,
		"shutdownBanks":         shutdownBanksStatus,
		"withdrawShutdownBanks": cr.withdrawShutdownBanks,
		"insertShutdownBanks":   cr.insertShutdownBanks,
	}
}
