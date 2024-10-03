package pwr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"worldofnuclear.com/internal/simworks"
)

func TestNewGenerator(t *testing.T) {
	g := NewGenerator("Test Generator", "A test generator", nil)
	s := simworks.NewSimulator("Test Simulator", "A test simulator")
	s.AddComponent(g)

	_, err := g.Update(s)
	assert.Error(t, err)
}

func TestGeneratorNormalOperation(t *testing.T) {

	assert.Fail(t, "Test not implemented")
}
