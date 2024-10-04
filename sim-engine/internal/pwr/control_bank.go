package pwr

type ControlBank struct {
	label    string
	numRods  int
	position int // in steps, from 0 to MaxWithdrawalSteps
}

func NewControlBank(label string, numRods int) *ControlBank {
	return &ControlBank{
		label:   label,
		numRods: numRods,
	}
}

func (cb *ControlBank) Status() map[string]interface{} {
	return map[string]interface{}{
		"label":    cb.label,
		"numRods":  cb.numRods,
		"position": cb.position,
	}
}

func (cb *ControlBank) Label() string {
	return cb.label
}

func (cb *ControlBank) NumRods() int {
	return cb.numRods
}

func (cb *ControlBank) Position() int {
	return cb.position
}

func (cb *ControlBank) LowerPosition(steps int) {
	cb.position = max(cb.position-steps, 0)
}

func (cb *ControlBank) RaisePosition(steps int) {
	cb.position = min(cb.position+steps, int(Config["control_rods"]["max_withdrawal_steps"]))
}

func (cb *ControlBank) IsFullyWithdrawn() bool {
	return cb.position == int(Config["control_rods"]["max_withdrawal_steps"])
}

func (cb *ControlBank) IsFullyInserted() bool {
	return cb.position == 0
}

func (cb *ControlBank) Scram() {
	cb.position = 0
}
