package pwr

import (
	"errors"
	"fmt"
	"math"

	"worldofnuclear.com/internal/simworks"
)

/*
Notes:
- Steam turbine is a component that converts steam energy into mechanical energy.
- The steam turbine is driven by the steam pressure and flow rate from the steam generator.
- The steam turbine is used to drive the generator.

In order to turn at just the right speed (RPMs), the generator has a throttle valve that
controls the flow of steam into the turbine. It has a governor system that controls the
throttle valve. For our purposes to keep things simple, we will assume the throttle just works
when the steam pressure is high enough.

There is a synchronization process to get the generator to match the frequency of the grid.
This is done by a synchro system that sends a signal to the turbine to slow down or speed up.
The RPMs need to start slightly above target speed. Then it looks at the generator
output frequency. Then the phase angle is matched between the generator output and the grid.
When frequency and phase align, the breaker is closed and the generator connects to the grid.

We will model this using feedback between the steam turbine and the generator. Internally,
we might introduce a Governor that manages the coordination once the operator initiates
the synchronization process (with the intent to connect to the grid).

- The inlet pressure is the pressure of the steam entering the turbine.
- The outlet pressure is the pressure of the steam exiting the turbine.
- The steam flow rate is the rate at which steam is flowing through the turbine. This impacts
the RPMs.
- Thermal power is how much heat energy is converted into mechanical energy. The leftover
heat energy is wasted, sent to the condenser. Some amount of waste is unavoidable due to the
laws of thermodynamics.
- Blade diameter is the diameter of the turbine blades. This impacts the RPMs.
- Efficiency is the ratio of the useful work output to the work input.
- Max RPM is the maximum RPMs the turbine can operate at before bad things happen.
- RPM is the current RPMs of the turbine.
*/
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
