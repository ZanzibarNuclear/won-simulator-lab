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
	Update(env *Environment, otherComponents []Component)
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

func FindBoiler(components []Component) *Boiler {
	for _, component := range components {
		if boiler, ok := component.(*Boiler); ok {
			return boiler
		}
	}
	return nil
}

func FindTurbine(components []Component) *Turbine {
	for _, component := range components {
		if turbine, ok := component.(*Turbine); ok {
			return turbine
		}
	}
	return nil
}


