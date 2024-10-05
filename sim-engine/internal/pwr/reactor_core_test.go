package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"worldofnuclear.com/internal/simworks"
)

func TestReactorCore_Init(t *testing.T) {
	pl := NewPrimaryLoop("Test Primary Loop", "A test primary loop")
	r := NewReactorCore("Test Reactor Core", "A test reactor core", pl)
	s := simworks.NewSimulator("Test Simulator", "A test simulator")
	s.AddComponent(r)

	status, err := r.Update(s)
	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "Test Reactor Core", r.Name())
	assert.Equal(t, "A test reactor core", r.Description())
}

func TestReactorCore_NormalOperation(t *testing.T) {
	assert.Fail(t, "Test not implemented")
}
