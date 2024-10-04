package pwr

import (
	"errors"
	"fmt"
	"math"

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
func NewCondenser(name string, description string, steamTurbine *SteamTurbine, secondaryLoop *SecondaryLoop) *Condenser {
	return &Condenser{
		BaseComponent: *simworks.NewBaseComponent(name, description),
		steamTurbine:  steamTurbine,
		secondaryLoop: secondaryLoop,
		surfaceArea:   Config["condenser"]["surface_area"],
		tubeMaterial:  "titanium",
		condenserType: "water-cooled",
		waterVelocity: 2.0,
	}
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

// CalculateCondenserProperties calculates steam condenser properties for a PWR plant
func (c *Condenser) Update(s *simworks.Simulator) (map[string]interface{}, error) {

	// no events for condenser; just reaction to steam turbine and tertiary loop (cooling water)

	if c.steamTurbine == nil {
		fmt.Println("Error: Steam Turbine not found")
		return nil, errors.New("steam turbine not found")
	}

	// Material thermal conductivities (W/m·K)
	tubeConductivity := map[string]float64{
		"titanium":        22,
		"stainless_steel": 16,
		"copper_nickel":   50,
	}

	// Calculate thermal efficiency (typical PWR values)
	thermalEfficiency := Config["condenser"]["thermal_efficiency"]

	// Calculate heat rejection to condenser
	// FIXME: get thermal power from steam turbine??
	heatRejectionMW := c.steamTurbine.ThermalPower() * (1 - thermalEfficiency)
	heatRejectionW := heatRejectionMW * 1e6

	// Calculate condenser pressure based on cooling type and temperature
	var ttd, condenserTemp float64
	if c.condenserType == "water-cooled" {
		ttd = 5.0                                      //
		condenserTemp = c.coolingWaterTempIn + ttd + 5 // Adding approach temperature
	} else { // air-cooled
		ttd = 15.0 // Larger TTD for air-cooled
		condenserTemp = c.coolingWaterTempIn + ttd + 10
	}

	// Calculate saturation pressure using simplified correlation
	c.condenserPressure = 3.17 * math.Exp(0.0724*condenserTemp)
	c.condenserTemperature = condenserTemp

	// Calculate steam flow rate
	c.steamFlowRate = heatRejectionW / Config["common"]["steam_latent_heat"]

	// Calculate cooling water flow rate for water-cooled condenser
	if c.condenserType == "water-cooled" {
		tempRise := 10.0 // Typical temperature rise in cooling water (°C)
		c.coolingWaterFlowRate = heatRejectionW / (Config["common"]["water_specific_heat"] * tempRise)

		// Calculate tube parameters
		tubeOD := 0.025         // m (typical 25mm outer diameter)
		tubeThickness := 0.0015 // m (typical 1.5mm thickness)
		tubeID := tubeOD - 2*tubeThickness

		// Calculate heat transfer coefficients
		k := tubeConductivity[c.tubeMaterial]
		reynolds := (c.waterVelocity * tubeID * Config["common"]["water_density"]) / 1e-3 // Using water viscosity ≈ 1e-3
		prandtl := 7.0                                                                    // Typical for water

		// Dittus-Boelter correlation for internal flow
		nusselt := 0.023 * math.Pow(reynolds, 0.8) * math.Pow(prandtl, 0.4)
		hWater := nusselt * 0.6 / tubeID // W/m²·K

		// Typical steam-side heat transfer coefficient
		hSteam := 10000.0 // W/m²·K

		// Overall heat transfer coefficient
		uOverall := 1 / (1/hSteam + tubeOD*math.Log(tubeOD/tubeID)/(2*k) + tubeOD/(tubeID*hWater))

		// Calculate required surface area
		lmtd := ((c.condenserTemperature - c.coolingWaterTempIn) -
			(c.condenserTemperature - (c.coolingWaterTempIn + tempRise))) /
			math.Log((c.condenserTemperature-c.coolingWaterTempIn)/
				(c.condenserTemperature-(c.coolingWaterTempIn+tempRise)))
		surfaceArea := heatRejectionW / (uOverall * lmtd)

		c.heatTransferCoefficient = uOverall
		c.surfaceArea = surfaceArea
	}
	return c.Status(), nil
}
