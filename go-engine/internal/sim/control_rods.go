package sim

type ControlRods struct {
	controlBanks          [7]*ControlBank
	shutdownBanks         [4]*ShutdownBank
	withdrawShutdownBanks bool // rods go up when true, down when false
}

func NewControlRods() *ControlRods {
	cr := &ControlRods{}

	// Initialize Control Banks
	for i := 0; i < 5; i++ {
		cr.controlBanks[i] = NewControlBank(4) // Assuming 4 rods per bank
	}
	cr.controlBanks[5] = NewControlBank(8)
	cr.controlBanks[6] = NewControlBank(9)

	// Initialize Shutdown Banks
	for i := 0; i < 4; i++ {
		cr.shutdownBanks[i] = NewShutdownBank(8) // Assuming 4 rods per bank
	}

	return cr
}

// bank 1 thru 7
func (cr *ControlRods) AdjustControlBankPosition(bank int, target int) {
	if bank < 1 || bank > 7 {
		// throw out of bounds error
		// panic("bank out of bounds")
		return
	}
	cr.controlBanks[bank-1].SetTarget(target)
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
	for _, bank := range cr.shutdownBanks {
		bank.Update()
	}
}

func (cr *ControlRods) Status() map[string]interface{} {
	controlBankPositions := make([]int, len(cr.controlBanks))
	for i, bank := range cr.controlBanks {
		controlBankPositions[i] = bank.Position()
	}

	shutdownBankPositions := make([]int, len(cr.shutdownBanks))
	for i, bank := range cr.shutdownBanks {
		shutdownBankPositions[i] = bank.Position()
	}

	return map[string]interface{}{
		"controlBanks":  controlBankPositions,
		"shutdownBanks": shutdownBankPositions,
	}
}
