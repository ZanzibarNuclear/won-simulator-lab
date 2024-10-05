package pwr

import (
	"errors"
	"fmt"

	"worldofnuclear.com/internal/simworks"
)

// CondenserProperties represents the calculated parameters for a PWR steam condenser
type Condenser struct {
	simworks.BaseComponent
	heatRejection           float64 // 2000 - 2500 MW (for 1000 MWe plant)
	condenserPressure       float64 // 0.003 - 0.008 MPa
	condenserTemperature    float64 // 30 - 45 °C
	steamFlowRate           float64 // 1500 - 2500 kg/s (for 1000 MWe plant)
	coolingWaterFlowRate    float64 // 40000 - 60000 kg/s
	heatTransferCoefficient float64 // 2000 - 4000 W/m²·K
	coolingWaterTempIn      float64 // 20 - 30 °C
	coolingWaterTempOut     float64 // rise of 8 - 12 °C
	waterVelocity           float64 // m/s
	surfaceArea             float64 // 30000 - 50000 m²
	tubeMaterial            string  // "titanium", "stainless_steel", "copper_nickel"
	condenserType           string  // "water-cooled", "air-cooled"
	steamTurbine            *SteamTurbine
	secondaryLoop           *SecondaryLoop
}

// NewCondenser creates a new Condenser instance
func NewCondenser(name string, description string, st *SteamTurbine, sl *SecondaryLoop) *Condenser {
	return &Condenser{
		BaseComponent: *simworks.NewBaseComponent(name, description),
		steamTurbine:  st,
		secondaryLoop: sl,
		surfaceArea:   float64(Config["condenser"]["surface_area"]),
		tubeMaterial:  "titanium",
		condenserType: "water-cooled",
		waterVelocity: 2.0, // m/s; expected?
	}
}

func (c *Condenser) SteamTurbine() *SteamTurbine {
	return c.steamTurbine
}

func (c *Condenser) HeatRejection() float64 {
	return c.heatRejection
}

func (c *Condenser) CondenserPressure() float64 {
	return c.condenserPressure
}

func (c *Condenser) CondenserTemperature() float64 {
	return c.condenserTemperature
}

func (c *Condenser) SteamFlowRate() float64 {
	return c.steamFlowRate
}

func (c *Condenser) CoolingWaterFlowRate() float64 {
	return c.coolingWaterFlowRate
}

func (c *Condenser) HeatTransferCoefficient() float64 {
	return c.heatTransferCoefficient
}

func (c *Condenser) CoolingWaterTempIn() float64 {
	return c.coolingWaterTempIn
}

func (c *Condenser) CoolingWaterTempOut() float64 {
	return c.coolingWaterTempOut
}

func (c *Condenser) WaterVelocity() float64 {
	return c.waterVelocity
}

func (c *Condenser) SurfaceArea() float64 {
	return c.surfaceArea
}

func (c *Condenser) TubeMaterial() string {
	return c.tubeMaterial
}

func (c *Condenser) CondenserType() string {
	return c.condenserType
}

func (c *Condenser) Status() map[string]interface{} {
	return map[string]interface{}{
		"about":                     c.BaseComponent.Status(),
		"heat_rejection":            c.HeatRejection(),
		"condenser_pressure":        c.CondenserPressure(),
		"condenser_temperature":     c.CondenserTemperature(),
		"steam_flow_rate":           c.SteamFlowRate(),
		"cooling_water_flow_rate":   c.CoolingWaterFlowRate(),
		"heat_transfer_coefficient": c.HeatTransferCoefficient(),
		"cooling_water_temp_in":     c.CoolingWaterTempIn(),
		"cooling_water_temp_out":    c.CoolingWaterTempOut(),
		"water_velocity":            c.WaterVelocity(),
		"surface_area":              c.SurfaceArea(),
	}
}

// PrintResults prints the calculated results in a formatted way
func (c *Condenser) Print() {
	fmt.Println("==> Steam Condenser")
	c.BaseComponent.Print()
	fmt.Printf("Heat Rejection: %.0f MW\n", c.HeatRejection())
	fmt.Printf("Condenser Type: %s\n", c.CondenserType())

	fmt.Println("\nOperating Conditions:")
	fmt.Printf("Condenser Pressure: %.2f kPa\n", c.CondenserPressure())
	fmt.Printf("Condenser Temperature: %.1f°C\n", c.CondenserTemperature())
	fmt.Printf("Steam Flow Rate: %.1f kg/s\n", c.SteamFlowRate())

	if c.CondenserType() == "water-cooled" {
		fmt.Println("\nCooling Water Parameters:")
		fmt.Printf("Cooling Water Flow Rate: %.1f kg/s\n", c.CoolingWaterFlowRate())
		fmt.Printf("Heat Transfer Coefficient: %.1f W/m²·K\n", c.HeatTransferCoefficient())
		fmt.Printf("Cooling Water Temp In: %.1f°C\n", c.CoolingWaterTempIn())
		fmt.Printf("Cooling Water Temp Out: %.1f°C\n", c.CoolingWaterTempOut())
		fmt.Printf("Water Velocity: %.1f m/s\n", c.WaterVelocity())
		fmt.Printf("Heat Transfer Surface Area: %.1f m²\n", c.SurfaceArea())
	} else {
		fmt.Println("\nAir-cooled Parameters:")
		fmt.Println("Add air-cooled information here...")
	}
}

func (c *Condenser) Update(s *simworks.Simulator) (map[string]interface{}, error) {

	if c.steamTurbine == nil || c.secondaryLoop == nil {
		return nil, errors.New("steam turbine and secondary loop are required")
	}

	efficiency := Config["steam_turbine"]["efficiency"]
	c.heatRejection = c.SteamTurbine().ThermalPower() * (1 - efficiency)

	c.coolingWaterTempIn = Config["condenser"]["cooling_water_temp_in"]
	c.coolingWaterTempOut = Config["condenser"]["cooling_water_temp_out"]

	return c.Status(), nil
}
