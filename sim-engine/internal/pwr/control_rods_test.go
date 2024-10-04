package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestControlRods_Init(t *testing.T) {
	cr := NewControlRods()

	// Test control banks
	assert.Equal(t, "MA1", cr.controlBanks[0].Label())
	assert.Equal(t, 4, cr.controlBanks[0].NumRods())
	assert.Equal(t, 0, cr.controlBanks[0].Position())

	assert.Equal(t, "MA2", cr.controlBanks[1].Label())
	assert.Equal(t, 4, cr.controlBanks[1].NumRods())
	assert.Equal(t, 0, cr.controlBanks[1].Position())

	assert.Equal(t, "MB1", cr.controlBanks[2].Label())
	assert.Equal(t, 4, cr.controlBanks[2].NumRods())
	assert.Equal(t, 0, cr.controlBanks[2].Position())

	assert.Equal(t, "MB2", cr.controlBanks[3].Label())
	assert.Equal(t, 4, cr.controlBanks[3].NumRods())
	assert.Equal(t, 0, cr.controlBanks[3].Position())

	// Test gray banks
	assert.Equal(t, "GR1", cr.grayBanks[0].Label())
	assert.Equal(t, 8, cr.grayBanks[0].NumRods())
	assert.Equal(t, 0, cr.grayBanks[0].Position())

	assert.Equal(t, "GR2", cr.grayBanks[1].Label())
	assert.Equal(t, 8, cr.grayBanks[1].NumRods())
	assert.Equal(t, 0, cr.grayBanks[1].Position())

	// Test shutdown banks
	assert.Equal(t, "SD1", cr.shutdownBanks[0].Label())
	assert.Equal(t, 8, cr.shutdownBanks[0].NumRods())
	assert.Equal(t, 0, cr.shutdownBanks[0].Position())

	assert.Equal(t, "SD2", cr.shutdownBanks[1].Label())
	assert.Equal(t, 8, cr.shutdownBanks[1].NumRods())
	assert.Equal(t, 0, cr.shutdownBanks[1].Position())

	assert.Equal(t, "SD3", cr.shutdownBanks[2].Label())
	assert.Equal(t, 8, cr.shutdownBanks[2].NumRods())
	assert.Equal(t, 0, cr.shutdownBanks[2].Position())

	assert.Equal(t, "SD4", cr.shutdownBanks[3].Label())
	assert.Equal(t, 8, cr.shutdownBanks[3].NumRods())
	assert.Equal(t, 0, cr.shutdownBanks[3].Position())
}

func TestControlRods_RaiseAndLowerControlBanks(t *testing.T) {
	cr := NewControlRods()

	var cnt int

	// raise up a bit
	targetPosition := 42
	for cnt = 0; cr.controlBanks[0].Position() < targetPosition && cnt < 300; cnt++ {
		cr.AdjustControlBanks(targetPosition)
	}
	assert.Equal(t, targetPosition, cr.controlBanks[0].Position())

	// up more
	targetPosition = 125
	for cnt = 0; cr.controlBanks[0].Position() < targetPosition && cnt < 300; cnt++ {
		cr.AdjustControlBanks(targetPosition)
	}
	assert.Equal(t, targetPosition, cr.controlBanks[0].Position())
	assert.Equal(t, cr.AverageControlRodExtraction(), 0.5)

	// down
	targetPosition = 75
	for cnt = 0; cr.controlBanks[0].Position() > targetPosition && cnt < 300; cnt++ {
		cr.AdjustControlBanks(targetPosition)
	}
	assert.Equal(t, targetPosition, cr.controlBanks[0].Position())

	// all the way down
	targetPosition = 0
	for cnt = 0; cr.controlBanks[0].Position() > targetPosition && cnt < 300; cnt++ {
		cr.AdjustControlBanks(targetPosition)
	}
	assert.Equal(t, targetPosition, cr.controlBanks[0].Position())
}

func TestControlRods_RaiseAndLowerGrayBanks(t *testing.T) {
	cr := NewControlRods()

	var cnt int

	// raise up a bit
	targetPosition := 42
	for cnt = 0; cr.grayBanks[0].Position() < targetPosition && cnt < 300; cnt++ {
		cr.AdjustGrayBanks(targetPosition)
	}
	assert.Equal(t, targetPosition, cr.grayBanks[0].Position())

	// up more
	targetPosition = 103
	for cnt = 0; cr.grayBanks[0].Position() < targetPosition && cnt < 300; cnt++ {
		cr.AdjustGrayBanks(targetPosition)
	}
	assert.Equal(t, targetPosition, cr.grayBanks[0].Position())

	// down
	targetPosition = 75
	for cnt = 0; cr.grayBanks[0].Position() > targetPosition && cnt < 300; cnt++ {
		cr.AdjustGrayBanks(targetPosition)
	}
	assert.Equal(t, targetPosition, cr.grayBanks[0].Position())

	// all the way down
	targetPosition = 0
	for cnt = 0; cr.grayBanks[0].Position() > targetPosition && cnt < 300; cnt++ {
		cr.AdjustGrayBanks(targetPosition)
	}
	assert.Equal(t, targetPosition, cr.grayBanks[0].Position())
}

func TestControlRods_WithdrawInsertShutdownBanks(t *testing.T) {
	cr := NewControlRods()

	var cnt int
	for cnt = 0; !cr.ShutdownBanksFullyWithdrawn() && cnt < 300; cnt++ {
		cr.WithdrawShutdownBanks()
	}

	assert.Greater(t, cnt, 50) // not too fast
	assert.True(t, cr.shutdownBanks[0].IsFullyWithdrawn())
	assert.True(t, cr.shutdownBanks[1].IsFullyWithdrawn())
	assert.True(t, cr.shutdownBanks[2].IsFullyWithdrawn())
	assert.True(t, cr.shutdownBanks[3].IsFullyWithdrawn())

	// now the other direction
	for cnt = 0; !cr.ShutdownBanksFullyInserted() && cnt < 300; cnt++ {
		cr.InsertShutdownBanks()
	}

	assert.Greater(t, cnt, 50) // not too fast
	assert.True(t, cr.shutdownBanks[0].IsFullyInserted())
	assert.True(t, cr.shutdownBanks[1].IsFullyInserted())
	assert.True(t, cr.shutdownBanks[2].IsFullyInserted())
	assert.True(t, cr.shutdownBanks[3].IsFullyInserted())
}

func TestControlRods_ScramControlRods(t *testing.T) {
	cr := NewControlRods()

	initialPosition := 200

	// Raise the position of control banks to a non-zero value
	for _, bank := range cr.controlBanks {
		bank.RaisePosition(initialPosition)
		assert.Equal(t, initialPosition, bank.Position())
	}

	// Raise the position of gray banks to a non-zero value
	for _, bank := range cr.grayBanks {
		bank.RaisePosition(initialPosition)
		assert.Equal(t, initialPosition, bank.Position())
	}

	// Raise the position of shutdown banks to a non-zero value
	for _, bank := range cr.shutdownBanks {
		bank.RaisePosition(initialPosition)
		assert.Equal(t, initialPosition, bank.Position())
	}

	// Test scram
	cr.Scram()
	for _, bank := range cr.controlBanks {
		assert.True(t, bank.IsFullyInserted())
	}
	for _, bank := range cr.grayBanks {
		assert.True(t, bank.IsFullyInserted())
	}
	for _, bank := range cr.shutdownBanks {
		assert.True(t, bank.IsFullyInserted())
	}
}
