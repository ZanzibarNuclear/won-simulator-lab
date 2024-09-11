package sim

import (
	"fmt"
	"math"
)

type SteamTurbine struct {
	BaseComponent
	rpm           int     // Revolutions per minute
	maxRPM        int     // Maximum RPM the turbine can handle
	efficiency    float64 // Turbine efficiency (0-1)
	steamPressure float64 // Current steam pressure from SteamGenerator (in Pascal)
}

func NewSteamTurbine(name string) *SteamTurbine {
	return &SteamTurbine{
		BaseComponent: BaseComponent{Name: name},
		rpm:           0,
		maxRPM:        TURBINE_MAX_RPM,
		efficiency:    0.9, // 90% efficiency, can be adjusted
		steamPressure: 0,
	}
}

func (st *SteamTurbine) Update(env *Environment, s *Simulation) {
	steamGen := s.FindSteamGenerator()
	if steamGen == nil {
		fmt.Println("Error: Steam Generator not found")
		return
	}

	// Update steam pressure based on SteamGenerator's output
	st.steamPressure = steamGen.steamFlowRate * 1000 // Simple conversion, adjust as needed

	// Calculate RPM based on steam pressure
	// This is a simplified calculation and should be replaced with a more accurate model
	targetRPM := int(st.steamPressure / 10000 * float64(st.maxRPM) * st.efficiency)
	
	// Gradually adjust RPM (turbines don't instantly change speed)
	rpmDiff := targetRPM - st.rpm
	st.rpm += int(float64(rpmDiff) * 0.1) // Adjust 10% of the difference

	// Ensure RPM stays within bounds
	st.rpm = int(math.Max(0, math.Min(float64(st.rpm), float64(st.maxRPM))))
}

func (st *SteamTurbine) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":          st.Name,
		"rpm":           st.rpm,
		"maxRPM":        st.maxRPM,
		"efficiency":    st.efficiency,
		"steamPressure": st.steamPressure,
	}
}

func (st *SteamTurbine) PrintStatus() {
	fmt.Printf("Steam Turbine: %s\n", st.Name)
	fmt.Printf("\tRPM: %d\n", st.rpm)
	fmt.Printf("\tMax RPM: %d\n", st.maxRPM)
	fmt.Printf("\tEfficiency: %.2f\n", st.efficiency)
	fmt.Printf("\tSteam Pressure: %.2f Pa\n", st.steamPressure)
}

func (st *SteamTurbine) Rpm() int {
	return st.rpm
}
