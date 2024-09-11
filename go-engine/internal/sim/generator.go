package sim

import (
	"fmt"
)

type Generator struct {
	BaseComponent
	rpm             float64
	electricalPower float64 // in megawatts (MW)
}

func NewGenerator(name string) *Generator {
	return &Generator{
		BaseComponent:   BaseComponent{Name: name},
		rpm:             0,
		electricalPower: 0,
	}
}

func (g *Generator) Update(env *Environment, s *Simulation) {
	turbine := s.FindSteamTurbine()
	if turbine == nil {
		fmt.Println("No turbine found")
		return
	}

	// Update RPM based on turbine's RPM
	g.rpm = float64(turbine.Rpm())

	// Simple calculation of electrical power based on RPM
	// This is a simplified model and should be replaced with a more accurate one
	g.electricalPower = g.rpm * 0.001 // Arbitrary scaling factor
}

func (g *Generator) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":            g.Name,
		"rpm":             g.rpm,
		"electricalPower": g.electricalPower,
	}
}

func (g *Generator) PrintStatus() {
	fmt.Printf("Generator: %s\n", g.Name)
	fmt.Printf("\tRPM: %.2f\n", g.rpm)
	fmt.Printf("\tElectrical Power: %.2f MW\n", g.electricalPower)
}

func (g *Generator) GetRPM() float64 {
	return g.rpm
}

func (g *Generator) GetElectricalPower() float64 {
	return g.electricalPower
}
