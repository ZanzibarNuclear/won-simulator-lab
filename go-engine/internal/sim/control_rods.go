package sim

const (
	MAX_WITHDRAWAL_STEPS = 200
	CONTROL_ROD_SPEED    = 40 // steps per minute
)

type ControlBank struct {
	position int // in steps
	numRods  int
}

func (cb *ControlBank) FullyWithdrawn() bool {
	return cb.position == MAX_WITHDRAWAL_STEPS
}

func (cb *ControlBank) FullyInserted() bool {
	return cb.position == 0
}

func (cb *ControlBank) MoveTowardTarget(target int) {
	if target > cb.position {
		cb.position += min(CONTROL_ROD_SPEED, target-cb.position)
	} else if target < cb.position {
		cb.position -= min(CONTROL_ROD_SPEED, cb.position-target)
	}
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
		for _, bank := range cr.shutdownBanks {
			if !bank.FullyWithdrawn() {
				bank.MoveTowardTarget(MAX_WITHDRAWAL_STEPS)
				break // Only start withdrawing one bank at a time
			}
		}
	}

	if cr.insertShutdownBanks && !cr.ShutdownBanksFullyInserted() {
		for i := len(cr.shutdownBanks) - 1; i >= 0; i-- {
			bank := cr.shutdownBanks[i]
			if !bank.FullyInserted() {
				bank.MoveTowardTarget(0)
				break
			}
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
		if !bank.FullyWithdrawn() {
			return false
		}
	}
	return true
}

func (cr *ControlRods) ShutdownBanksFullyInserted() bool {
	for _, bank := range cr.shutdownBanks {
		if !bank.FullyInserted() {
			return false
		}
	}
	return true
}

func (cr *ControlRods) Status() map[string]interface{} {
	controlBanksStatus := make([]map[string]interface{}, len(cr.controlBanks))
	for i, bank := range cr.controlBanks {
		controlBanksStatus[i] = map[string]interface{}{
			"controlBankNum": i,
			"position":       bank.position,
		}
	}

	shutdownBanksStatus := make([]map[string]interface{}, len(cr.shutdownBanks))
	for i, bank := range cr.shutdownBanks {
		shutdownBanksStatus[i] = map[string]interface{}{
			"shutdownBankNum": i,
			"position":        bank.position,
		}
	}

	return map[string]interface{}{
		"controlBanks":          controlBanksStatus,
		"shutdownBanks":         shutdownBanksStatus,
		"withdrawShutdownBanks": cr.withdrawShutdownBanks,
		"insertShutdownBanks":   cr.insertShutdownBanks,
	}
}
