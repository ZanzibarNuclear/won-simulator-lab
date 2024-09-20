package sim

const (
	MAX_WITHDRAWAL_STEPS = 250
	WITHDRAWAL_RATE      = 50 // steps per minute
)

type ControlBank struct {
	label    string
	numRods  int
	position int // in steps, from 0 to MaxWithdrawalSteps
	target   int // target position
}

func NewControlBank(label string, numRods int) *ControlBank {
	return &ControlBank{
		label:    label,
		numRods:  numRods,
		position: 0,
		target:   0,
	}
}

func (cb *ControlBank) Label() string {
	return cb.label
}

func (cb *ControlBank) Position() int {
	return cb.position
}

func (cb *ControlBank) NumRods() int {
	return cb.numRods
}

func (cb *ControlBank) Target() int {
	return cb.target
}

func (cb *ControlBank) SetTarget(target int) {
	if target < 0 {
		cb.target = 0
	} else if target > MAX_WITHDRAWAL_STEPS {
		cb.target = MAX_WITHDRAWAL_STEPS
	} else {
		cb.target = target
	}
}

func (cb *ControlBank) Update() {
	if cb.position < cb.target {
		cb.position = min(cb.position+WITHDRAWAL_RATE, cb.target)
	} else if cb.position > cb.target {
		cb.position = max(cb.position-WITHDRAWAL_RATE, cb.target)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (cb *ControlBank) IsFullyWithdrawn() bool {
	return cb.position == MAX_WITHDRAWAL_STEPS
}

func (cb *ControlBank) IsFullyInserted() bool {
	return cb.position == 0
}

func (cb *ControlBank) Scram() {
	cb.position = 0
	cb.target = 0
}

type ShutdownBank struct {
	ControlBank
}

func NewShutdownBank(label string, numRods int) *ShutdownBank {
	return &ShutdownBank{
		ControlBank: ControlBank{
			label:    label,
			numRods:  numRods,
			position: 0,
			target:   0,
		},
	}
}

func (sb *ShutdownBank) Withdraw() {
	sb.SetTarget(MAX_WITHDRAWAL_STEPS)
}

func (sb *ShutdownBank) Insert() {
	sb.SetTarget(0)
}
