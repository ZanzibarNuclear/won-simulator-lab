package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestControlBank_Init(t *testing.T) {
	cb := NewControlBank("TestBank", 10)

	assert.Equal(t, "TestBank", cb.Label())
	assert.Equal(t, 10, cb.NumRods())
	assert.Equal(t, 0, cb.Position())
}

func TestControlBank_RaiseAndLower(t *testing.T) {
	cb := NewControlBank("TestBank", 10)

	// Test raising the position
	cb.RaisePosition(5)
	assert.Equal(t, 5, cb.Position())
	assert.False(t, cb.IsFullyWithdrawn())
	assert.False(t, cb.IsFullyInserted())

	// Test raising the position beyond max withdrawal steps
	cb.RaisePosition(int(Config["control_rods"]["max_withdrawal_steps"]))
	assert.Equal(t, int(Config["control_rods"]["max_withdrawal_steps"]), cb.Position())
	assert.True(t, cb.IsFullyWithdrawn())
	assert.False(t, cb.IsFullyInserted())

	// Test lowering the position
	cb.LowerPosition(3)
	assert.Equal(t, int(Config["control_rods"]["max_withdrawal_steps"])-3, cb.Position())
	assert.False(t, cb.IsFullyWithdrawn())
	assert.False(t, cb.IsFullyInserted())

	// Test lowering the position below 0
	cb.LowerPosition(int(Config["control_rods"]["max_withdrawal_steps"]))
	assert.Equal(t, 0, cb.Position())
	assert.True(t, cb.IsFullyInserted())
	assert.False(t, cb.IsFullyWithdrawn())
}

func TestControlBank_Scram(t *testing.T) {
	cb := NewControlBank("TestBank", 10)

	// Raise the position to a non-zero value
	cb.RaisePosition(5)
	assert.Equal(t, 5, cb.Position())

	// Test scram
	cb.Scram()
	assert.True(t, cb.IsFullyInserted())
}
