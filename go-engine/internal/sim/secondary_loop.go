package sim

import (
	"fmt"
)

type SecondaryLoop struct {
	BaseComponent
	waterFlowRate float64 // in m³/s
	pumpPressure  float64 // in Pascal
}

func NewSecondaryLoop(name string) *SecondaryLoop {
	return &SecondaryLoop{
		BaseComponent: BaseComponent{Name: name},
		waterFlowRate: 100.0, // Initial value, can be adjusted as needed
		pumpPressure:  1000000.0, // 1 MPa, initial value
	}
}

func (sl *SecondaryLoop) Update(env *Environment, s *Simulation) {
	// Simplified update logic
	// In a real scenario, this would involve complex thermodynamics calculations
	
	// For example, we could adjust the water flow rate based on the turbine's RPM
	turbine := s.FindSteamTurbine()
	if turbine != nil {
		sl.waterFlowRate = float64(turbine.Rpm()) / TURBINE_MAX_RPM * 150.0 // Max flow rate of 150 m³/s
	}

	// Adjust pump pressure based on water flow rate (simplified relationship)
	sl.pumpPressure = 800000.0 + (sl.waterFlowRate * 5000.0) // Base pressure + dynamic pressure
}

func (sl *SecondaryLoop) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":          sl.Name,
		"waterFlowRate": sl.waterFlowRate,
		"pumpPressure":  sl.pumpPressure,
	}
}

func (sl *SecondaryLoop) PrintStatus() {
	fmt.Printf("Secondary Loop: %s\n", sl.Name)
	fmt.Printf("\tWater Flow Rate: %.2f m³/s\n", sl.waterFlowRate)
	fmt.Printf("\tPump Pressure: %.2f Pa\n", sl.pumpPressure)
}
