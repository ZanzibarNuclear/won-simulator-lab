package components

import (
	"won/sim-lab/go-engine/common"
)

type Component interface {
	Update(env *common.Environment, otherComponents []Component)
}
