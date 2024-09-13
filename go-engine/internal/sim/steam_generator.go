package sim

import (
	"fmt"
	"math"
)

type SteamGenerator struct {
	BaseComponent
	primaryInletTemp    float64 // Temperature of water coming from reactor core (°C)
	primaryOutletTemp   float64 // Temperature of water returning to reactor core (°C)
	secondaryInletTemp  float64 // Temperature of water from secondary loop (°C)
	secondaryOutletTemp float64 // Temperature of steam to secondary loop (°C)
	heatTransferRate    float64 // Rate of heat transfer from primary to secondary loop (MW)
	steamFlowRate       float64 // Rate of steam production (kg/s)
}

func NewSteamGenerator(name string) *SteamGenerator {
	return &SteamGenerator{
		BaseComponent:       BaseComponent{Name: name},
		primaryInletTemp:    320.0, // Initial values, can be adjusted as needed
		primaryOutletTemp:   280.0,
		secondaryInletTemp:  220.0,
		secondaryOutletTemp: 280.0,
		heatTransferRate:    1000.0, // 1000 MW, for example
		steamFlowRate:       500.0,  // 500 kg/s, for example
	}
}

func (sg *SteamGenerator) Update(env *Environment, s *Simulation) {
	reactorCore := s.FindReactorCore()
	secondaryLoop := s.FindSecondaryLoop()

	if reactorCore == nil || secondaryLoop == nil {
		fmt.Println("Error: Reactor Core or Secondary Loop not found")
		return
	}

	// Update primary inlet temperature based on reactor core heat
	sg.primaryInletTemp = math.Min(reactorCore.temperature, 350) // Max temp 350°C

	// Calculate heat transfer
	sg.heatTransferRate = reactorCore.GetHeatEnergyRate() * 0.95 // Assume 95% efficiency

	// Update temperatures
	tempDiff := sg.primaryInletTemp - sg.secondaryInletTemp
	sg.primaryOutletTemp = sg.primaryInletTemp - (tempDiff * 0.2)
	sg.secondaryOutletTemp = sg.secondaryInletTemp + (tempDiff * 0.8)

	// Calculate steam flow rate based on heat transfer
	// This is a simplified calculation and should be replaced with a more accurate model
	sg.steamFlowRate = sg.heatTransferRate * 0.5 // Arbitrary factor

	// Update secondary loop water flow rate
	secondaryLoop.feedwaterFlowRate = sg.steamFlowRate / 1000 // Convert kg/s to m³/s (assuming water density of 1000 kg/m³)
}

func (sg *SteamGenerator) Status() map[string]interface{} {
	return map[string]interface{}{
		"name":                sg.Name,
		"primaryInletTemp":    sg.primaryInletTemp,
		"primaryOutletTemp":   sg.primaryOutletTemp,
		"secondaryInletTemp":  sg.secondaryInletTemp,
		"secondaryOutletTemp": sg.secondaryOutletTemp,
		"heatTransferRate":    sg.heatTransferRate,
		"steamFlowRate":       sg.steamFlowRate,
	}
}

func (sg *SteamGenerator) PrintStatus() {
	fmt.Printf("Steam Generator: %s\n", sg.Name)
	fmt.Printf("\tPrimary Inlet Temperature: %.2f °C\n", sg.primaryInletTemp)
	fmt.Printf("\tPrimary Outlet Temperature: %.2f °C\n", sg.primaryOutletTemp)
	fmt.Printf("\tSecondary Inlet Temperature: %.2f °C\n", sg.secondaryInletTemp)
	fmt.Printf("\tSecondary Outlet Temperature: %.2f °C\n", sg.secondaryOutletTemp)
	fmt.Printf("\tHeat Transfer Rate: %.2f MW\n", sg.heatTransferRate)
	fmt.Printf("\tSteam Flow Rate: %.2f kg/s\n", sg.steamFlowRate)
}
