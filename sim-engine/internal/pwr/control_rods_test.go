package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewControlRods(t *testing.T) {
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
