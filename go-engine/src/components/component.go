// Copyright (c) 2024 Nuclear Ambitions LLC. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
//
// Package components provides a common interface for all components in a
// power plant, as well as an implementation of each specific component.
// While each component has its own characteristics, they share
// common properties and interactions with the simulation.
package components

import (
	"won/sim-lab/go-engine/common"
)

type Component interface {
	Update(env *common.Environment, otherComponents []Component)
}
