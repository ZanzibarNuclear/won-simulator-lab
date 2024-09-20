package sim

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

type ControlRods struct {
	controlBanks          [4]*ControlBank  // full-strength absorbers; for power control during operation
	grayBanks             [2]*ControlBank  // lower neutron absorption; for load following and fine reactivity control during operation
	shutdownBanks         [4]*ShutdownBank // full-strength absorbers; for shutting down the core rapidly
	withdrawShutdownBanks bool             // rods go up when true, down when false
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

	cr.shutdownBanks[0] = NewShutdownBank("SD1", 8)
	cr.shutdownBanks[1] = NewShutdownBank("SD2", 8)
	cr.shutdownBanks[2] = NewShutdownBank("SD3", 8)
	cr.shutdownBanks[3] = NewShutdownBank("SD4", 8)

	return cr
}

func (cr *ControlRods) InitiateShutdownBankWithdrawal() {
	cr.withdrawShutdownBanks = true
}

func (cr *ControlRods) InitiateShutdownBankInsertion() {
	cr.withdrawShutdownBanks = false
}

func (cr *ControlRods) ShutdownBanksFullyWithdrawn() bool {
	for _, bank := range cr.shutdownBanks {
		if !bank.IsFullyWithdrawn() {
			return false
		}
	}
	return true
}

func (cr *ControlRods) ShutdownBanksFullyInserted() bool {
	for _, bank := range cr.shutdownBanks {
		if !bank.IsFullyInserted() {
			return false
		}
	}
	return true
}

// bank 1 thru 7
func (cr *ControlRods) AdjustControlBankPosition(bank int, target int) {
	if bank < 1 || bank > 4 {
		fmt.Printf("Tried to move control bank %d, which is out of bounds/n", bank)
		return
	}
	cr.controlBanks[bank-1].SetTarget(target)
}

func (cr *ControlRods) AdjustGrayBankPosition(bank int, target int) {
	if bank < 1 || bank > 2 {
		fmt.Printf("Tried to move gray bank %d, which is out of bounds/n", bank)
		return
	}
	cr.grayBanks[bank-1].SetTarget(target)
}

func (cr *ControlRods) AverageControlRodExtraction() float64 {
	totalSteps := 0

	// never mind gray banks for now
	maxSteps := MAX_WITHDRAWAL_STEPS * len(cr.controlBanks)

	for _, bank := range cr.controlBanks {
		totalSteps += bank.Position()
	}

	if maxSteps == 0 {
		return 0 // avoid divide by zero; should not be possible though
	}

	return float64(totalSteps) / float64(maxSteps) * 100
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
	cr.withdrawShutdownBanks = false
}

func (cr *ControlRods) Update() {
	if cr.withdrawShutdownBanks {
		for _, bank := range cr.shutdownBanks {
			if !bank.IsFullyWithdrawn() {
				if bank.Target() == 0 {
					bank.Withdraw()
				}
				break
			}
		}
	} else {
		for _, bank := range cr.shutdownBanks {
			if !bank.IsFullyInserted() {
				if bank.Target() != 0 {
					bank.Insert()
				}
				break
			}
		}
	}

	for _, bank := range cr.controlBanks {
		bank.Update()
	}
	for _, bank := range cr.grayBanks {
		bank.Update()
	}
	for _, bank := range cr.shutdownBanks {
		bank.Update()
	}
}

func (cr *ControlRods) Status() map[string]interface{} {
	controlBankPositions := make([]int, len(cr.controlBanks))
	for i, bank := range cr.controlBanks {
		controlBankPositions[i] = bank.Position()
	}

	grayBankPositions := make([]int, len(cr.grayBanks))
	for i, bank := range cr.grayBanks {
		grayBankPositions[i] = bank.Position()
	}

	shutdownBankPositions := make([]int, len(cr.shutdownBanks))
	for i, bank := range cr.shutdownBanks {
		shutdownBankPositions[i] = bank.Position()
	}

	return map[string]interface{}{
		"controlBanks":  controlBankPositions,
		"grayBanks":     grayBankPositions,
		"shutdownBanks": shutdownBankPositions,
	}
}
