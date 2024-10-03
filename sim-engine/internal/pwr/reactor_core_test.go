package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"worldofnuclear.com/internal/simworks"
)

func TestNewReactorCore(t *testing.T) {
	r := NewReactorCore("Test Reactor Core", "A test reactor core")
	s := simworks.NewSimulator("Test Simulator", "A test simulator")
	s.AddComponent(r)

	status, err := r.Update(s)
	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "Test Reactor Core", r.Name())
	assert.Equal(t, "A test reactor core", r.Description())
}

func TestReactorCoreNormalOperation(t *testing.T) {
	assert.Fail(t, "Test not implemented")
}
