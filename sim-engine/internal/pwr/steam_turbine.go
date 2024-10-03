package pwr

import (
	"errors"
	"fmt"
	"math"

	"worldofnuclear.com/internal/simworks"
)

type SteamTurbine struct {
	simworks.BaseComponent
	inletPressure  float64 // MPa; check secondary loop
	outletPressure float64 // MPa; check secondary loop
	steamFlowRate  float64 // kg/s; check secondary loop
	rpm            int     // Revolutions per minute
	thermalPower   float64 // MW
	secondaryLoop  *SecondaryLoop
}

func NewSteamTurbine(name string, description string, secondaryLoop *SecondaryLoop) *SteamTurbine {
	return &SteamTurbine{
		BaseComponent: *simworks.NewBaseComponent(name, description),
		secondaryLoop: secondaryLoop,
	}
}

func (st *SteamTurbine) InletPressure() float64 {
	return st.inletPressure
}

func (st *SteamTurbine) OutletPressure() float64 {
	return st.outletPressure
}

func (st *SteamTurbine) SteamFlowRate() float64 {
	return st.steamFlowRate
}

func (st *SteamTurbine) BladeDiameter() float64 {
	return Config["steam_turbine"]["blade_diameter"]
}

func (st *SteamTurbine) Efficiency() float64 {
	return Config["steam_turbine"]["efficiency"]
}

func (st *SteamTurbine) MaxRPM() int {
	return int(Config["steam_turbine"]["max_rpm"])
}

func (st *SteamTurbine) Rpm() int {
	return st.rpm
}

func (st *SteamTurbine) ThermalPower() float64 {
	return st.thermalPower
}

func (st *SteamTurbine) Status() map[string]interface{} {
	return map[string]interface{}{
		"about":          st.BaseComponent.Status(),
		"inletPressure":  st.InletPressure(),
		"outletPressure": st.OutletPressure(),
		"steamFlowRate":  st.SteamFlowRate(),
		"rpm":            st.Rpm(),
		"bladeDiameter":  st.BladeDiameter(),
		"efficiency":     st.Efficiency(),
		"maxRPM":         st.MaxRPM(),
		"thermalPower":   st.ThermalPower(),
	}
}

func (st *SteamTurbine) Print() {
	fmt.Printf("=> Steam Turbine\n")
	st.BaseComponent.Print()
	fmt.Printf("Inlet pressure: %.2f MPa\n", st.InletPressure())
	fmt.Printf("Outlet pressure: %.2f MPa\n", st.OutletPressure())
	fmt.Printf("Steam flow rate: %.2f kg/s\n", st.SteamFlowRate())
	fmt.Printf("RPM: %d\n", st.Rpm())
	fmt.Printf("Blade diameter: %.2f m\n", st.BladeDiameter())
	fmt.Printf("Efficiency: %.2f\n", st.Efficiency())
	fmt.Printf("Max RPM: %d\n", st.MaxRPM())
	fmt.Printf("Thermal Power: %.2f MW\n", st.ThermalPower())
}

func (st *SteamTurbine) Update(s *simworks.Simulator) (map[string]interface{}, error) {

	// no direct events; just reaction to secondary loop values

	if st.secondaryLoop == nil {
		fmt.Println("Error: Secondary Loop not found")
		return nil, errors.New("secondary loop not found")
	}

	// Update steam pressure based on SteamGenerator's output
	st.inletPressure = st.secondaryLoop.steamPressure
	// st.steamFlowRate = st.secondaryLoop.steamFlowRate FIXME: add steam flow rate to secondary loop

	// st.outletPressure = st.secondaryLoop.steamPressure - (st.secondaryLoop.steamFlowRate * 1000) // FIXME: use a more accurate model

	st.rpm = st.CalculateRPM()
	st.thermalPower = 150.0 // FIXME: derive from turbine movement?? or use typical PWR values
	return st.Status(), nil
}

func (st *SteamTurbine) CalculateRPM() int {
	// Constants
	atmosphericPressure := Config["common"]["atmospheric_pressure"]
	gravity := Config["common"]["gravity"]
	efficiency := Config["steam_turbine"]["efficiency"]
	bladeDiameter := Config["steam_turbine"]["blade_diameter"]

	// Convert pressures to absolute
	inletPressureAbs := st.inletPressure + atmosphericPressure
	outletPressureAbs := st.outletPressure + atmosphericPressure

	// Calculate pressure ratio
	pressureRatio := inletPressureAbs / outletPressureAbs

	// Calculate blade tip speed using simplified Euler turbine equation principles
	bladeTipSpeed := math.Sqrt(
		2 * gravity * efficiency *
			(math.Pow(pressureRatio, 0.286) - 1) *
			(st.steamFlowRate),
	)

	// Calculate RPM using blade tip speed and diameter
	rpm := (bladeTipSpeed * 60) / (math.Pi * bladeDiameter)

	return int(math.Round(rpm))
}
