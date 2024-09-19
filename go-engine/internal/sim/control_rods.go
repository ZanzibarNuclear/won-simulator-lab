package sim

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

// bank 1 thru 7
func (cr *ControlRods) AdjustControlBankPosition(bank int, target int) {
	if bank < 1 || bank > 4 {
		// throw out of bounds error
		// panic("bank out of bounds")
		return
	}
	cr.controlBanks[bank-1].SetTarget(target)
}

func (cr *ControlRods) AdjustGrayBankPosition(bank int, target int) {
	if bank < 1 || bank > 2 {
		// throw out of bounds error
		// panic("bank out of bounds")
		return
	}
	cr.grayBanks[bank-1].SetTarget(target)
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
