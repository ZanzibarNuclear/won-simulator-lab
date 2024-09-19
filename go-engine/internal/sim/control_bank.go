package sim

const (
	MaxWithdrawalSteps = 200
)

type ControlBank struct {
	position int // in steps, from 0 to MaxWithdrawalSteps
	numRods  int
	target   int // target position
}

func NewControlBank(numRods int) *ControlBank {
	return &ControlBank{
		position: 0,
		numRods:  numRods,
		target:   0,
	}
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
	} else if target > MaxWithdrawalSteps {
		cb.target = MaxWithdrawalSteps
	} else {
		cb.target = target
	}
}

func (cb *ControlBank) Update() {
	if cb.position < cb.target {
		cb.position = min(cb.position+40, cb.target)
	} else if cb.position > cb.target {
		cb.position = max(cb.position-40, cb.target)
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
	return cb.position == MaxWithdrawalSteps
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

func NewShutdownBank(numRods int) *ShutdownBank {
	return &ShutdownBank{
		ControlBank: ControlBank{
			numRods:  numRods,
			position: 0,
			target:   0,
		},
	}
}

func (sb *ShutdownBank) Withdraw() {
	sb.SetTarget(MaxWithdrawalSteps)
}

func (sb *ShutdownBank) Insert() {
	sb.SetTarget(0)
}
