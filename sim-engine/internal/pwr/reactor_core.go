package pwr

import (
	"fmt"

	"worldofnuclear.com/internal/simworks"
)

type ReactorCore struct {
	simworks.BaseComponent
	fuelAge               int     // in minutes
	reactivity            float64 // negative means subcritical, 0 means critical, positive means supercritical
	neutronFlux           float64 // in neutrons per second
	temperature           float64 // in degrees Celsius
	heatEnergyRate        float64 // in MW
	controlRods           *ControlRods
	primaryLoop           *PrimaryLoop
	withdrawShutdownBanks bool
}

func NewReactorCore(name string, description string, primaryLoop *PrimaryLoop) *ReactorCore {
	return &ReactorCore{
		BaseComponent: *simworks.NewBaseComponent(name, description),
		primaryLoop:   primaryLoop,
	}
}

func (r *ReactorCore) FuelAge() int {
	return r.fuelAge
}

func (r *ReactorCore) SetFuelAge(age int) {
	r.fuelAge = age
}

func (r *ReactorCore) Reactivity() float64 {
	return r.reactivity
}

func (r *ReactorCore) SetReactivity(reactivity float64) {
	r.reactivity = reactivity
}

func (r *ReactorCore) NeutronFlux() float64 {
	return r.neutronFlux
}

func (r *ReactorCore) SetNeutronFlux(neutronFlux float64) {
	r.neutronFlux = neutronFlux
}

func (r *ReactorCore) Temperature() float64 {
	return r.temperature
}

func (r *ReactorCore) SetTemperature(temperature float64) {
	r.temperature = temperature
}

func (r *ReactorCore) HeatEnergyRate() float64 {
	return r.heatEnergyRate
}

func (r *ReactorCore) SetHeatEnergyRate(heatEnergyRate float64) {
	r.heatEnergyRate = heatEnergyRate
}

// TODO: add func to attach control rods

func (r *ReactorCore) Status() map[string]interface{} {
	return map[string]interface{}{
		"about": r.BaseComponent.Status(),
	}
}

func (r *ReactorCore) Print() {
	fmt.Printf("=> Reactor Core\n")
	r.BaseComponent.Print()
}

func (r *ReactorCore) Update(s *simworks.Simulator) (map[string]interface{}, error) {
	return r.Status(), nil
}
