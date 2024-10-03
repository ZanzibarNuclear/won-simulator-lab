package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewControlBank(t *testing.T) {
	cb := NewControlBank("TestBank", 10)

	assert.Equal(t, "TestBank", cb.Label())
	assert.Equal(t, 10, cb.NumRods())
	assert.Equal(t, 0, cb.Position())
	assert.Equal(t, 0, cb.Target())
}
