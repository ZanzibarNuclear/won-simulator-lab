package sim

const (
	MAX_STEPS = 200
	WITHDRAWAL_DURATION = 5 * 60 // 5 minutes in seconds
)

type ControlBank struct {
	steps   int
	numRods int
}

type ShutdownBank struct {
	insertion           float64
	withdrawalStartTime int64
	withdrawalEndTime   int64
	isWithdrawn         bool
	numRods             int
}

type ControlRods struct {
	controlBanks  [7]ControlBank
	shutdownBanks [4]*ShutdownBank
}

func NewControlRods() *ControlRods {
	cr := &ControlRods{}

	// Initialize control banks
	for i := 0; i < 5; i++ {
		cr.controlBanks[i] = ControlBank{steps: 0, numRods: 4}
	}
	cr.controlBanks[5] = ControlBank{steps: 0, numRods: 8}
	cr.controlBanks[6] = ControlBank{steps: 0, numRods: 9}

	// Initialize shutdown banks
	for i := 0; i < 4; i++ {
		cr.shutdownBanks[i] = &ShutdownBank{
			insertion:           1.0, // Fully inserted
			withdrawalStartTime: -1,
			withdrawalEndTime:   -1,
			isWithdrawn:         false,
			numRods:             8,
		}
	}

	return cr
}

func (cr *ControlRods) SetControlBankSteps(bankIndex, steps int) {
	if bankIndex < 0 || bankIndex >= len(cr.controlBanks) {
		return
	}
	cr.controlBanks[bankIndex].steps = max(0, min(steps, MAX_STEPS))
}

func (cr *ControlRods) InitiateShutdownBankWithdrawal(currentTime int64) {
	for _, bank := range cr.shutdownBanks {
		if !bank.isWithdrawn && bank.withdrawalStartTime == -1 {
			bank.withdrawalStartTime = currentTime
			bank.withdrawalEndTime = currentTime + WITHDRAWAL_DURATION
			break // Only start withdrawing one bank at a time
		}
	}
}

func (cr *ControlRods) UpdateShutdownBanks(currentTime int64) {
	for _, bank := range cr.shutdownBanks {
		if !bank.isWithdrawn && bank.withdrawalStartTime != -1 {
			if currentTime >= bank.withdrawalEndTime {
				bank.insertion = 0
				bank.isWithdrawn = true
			} else {
				progress := float64(currentTime - bank.withdrawalStartTime) / float64(bank.withdrawalEndTime - bank.withdrawalStartTime)
				bank.insertion = 1.0 - progress
			}
		}
	}
}

func (cr *ControlRods) CalculateAverageInsertion() float64 {
	totalSteps := 0
	totalRods := 0

	for _, bank := range cr.controlBanks {
		totalSteps += (MAX_STEPS - bank.steps) * bank.numRods
		totalRods += bank.numRods
	}

	for _, bank := range cr.shutdownBanks {
		totalSteps += int(bank.insertion * float64(MAX_STEPS) * float64(bank.numRods))
		totalRods += bank.numRods
	}

	return float64(totalSteps) / float64(totalRods*MAX_STEPS)
}

func (cr *ControlRods) ShutdownBanksStatus() []map[string]interface{} {
	status := make([]map[string]interface{}, len(cr.shutdownBanks))
	for i, bank := range cr.shutdownBanks {
		status[i] = map[string]interface{}{
			"insertion":   bank.insertion,
			"isWithdrawn": bank.isWithdrawn,
			"progress":    1.0 - bank.insertion,
		}
	}
	return status
}
