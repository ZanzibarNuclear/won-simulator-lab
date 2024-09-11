// Copyright (c) 2024 Nuclear Ambitions LLC. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
//
// Package components provides a common interface for all components in a
// power plant, as well as an implementation of each specific component.
// While each component has its own characteristics, they share
// common properties and interactions with the simulation.
package sim

type Component interface {
	Update(env *Environment, s *Simulation)
	PrintStatus()
	Status() map[string]interface{}
	GetName() string
}

type BaseComponent struct {
	Name string
}

func (b *BaseComponent) GetName() string {
	return b.Name
}

// ComponentType represents the type of a component in the power plant simulation.
type ComponentType int

const (
	Core ComponentType = iota
	Pump
	Condenser
	SteamTurbine
	Generator
)

// String returns the string representation of the ComponentType.
func (ct ComponentType) String() string {
	return [...]string{"Reactor Core", "Pump", "Condenser", "Steam Turbine", "Generator"}[ct]
}
