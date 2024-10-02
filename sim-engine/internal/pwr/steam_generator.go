package pwr

import (
	"errors"
	"fmt"

	"worldofnuclear.com/internal/simworks"
)

type SteamGenerator struct {
	simworks.BaseComponent
	primaryInletTemp    float64 // Temperature of hot leg
	primaryOutletTemp   float64 // Temperature of cold leg (simplification)
	secondaryInletTemp  float64 // Temperature of feedwater
	secondaryOutletTemp float64 // Temperature of steam out
	heatTransferRate    float64 // Rate of heat energy transfer to generate steam (MW)
	steamFlowRate       float64 // Rate of steam production (kg/s)
	primaryLoop         *PrimaryLoop
	secondaryLoop       *SecondaryLoop
}

func NewSteamGenerator(name string, description string, primaryLoop *PrimaryLoop, secondaryLoop *SecondaryLoop) *SteamGenerator {
	return &SteamGenerator{
		BaseComponent:       *simworks.NewBaseComponent(name, description),
		primaryLoop:         primaryLoop,
		secondaryLoop:       secondaryLoop,
		primaryInletTemp:    Config["common"]["room_temperature"],
		primaryOutletTemp:   Config["common"]["room_temperature"],
		secondaryInletTemp:  Config["common"]["room_temperature"],
		secondaryOutletTemp: Config["common"]["room_temperature"],
	}
}

func (sg *SteamGenerator) PrimaryInletTemp() float64 {
	return sg.primaryInletTemp
}

func (sg *SteamGenerator) PrimaryOutletTemp() float64 {
	return sg.primaryOutletTemp
}

func (sg *SteamGenerator) SecondaryInletTemp() float64 {
	return sg.secondaryInletTemp
}

func (sg *SteamGenerator) SecondaryOutletTemp() float64 {
	return sg.secondaryOutletTemp
}

func (sg *SteamGenerator) HeatTransferRate() float64 {
	return sg.heatTransferRate
}

func (sg *SteamGenerator) SteamFlowRate() float64 {
	return sg.steamFlowRate
}

func (sg *SteamGenerator) Status() map[string]interface{} {
	return map[string]interface{}{
		"about":               sg.BaseComponent.Status(),
		"primaryInletTemp":    sg.PrimaryInletTemp(),
		"primaryOutletTemp":   sg.PrimaryOutletTemp(),
		"secondaryInletTemp":  sg.SecondaryInletTemp(),
		"secondaryOutletTemp": sg.SecondaryOutletTemp(),
		"heatTransferRate":    sg.HeatTransferRate(),
		"steamFlowRate":       sg.SteamFlowRate(),
	}
}

func (sg *SteamGenerator) Print() {
	fmt.Printf("=> Steam Generator\n")
	sg.BaseComponent.Print()
	fmt.Printf("\tPrimary Inlet Temperature: %.2f °C\n", sg.PrimaryInletTemp())
	fmt.Printf("\tPrimary Outlet Temperature: %.2f °C\n", sg.PrimaryOutletTemp())
	fmt.Printf("\tSecondary Inlet Temperature: %.2f °C\n", sg.SecondaryInletTemp())
	fmt.Printf("\tSecondary Outlet Temperature: %.2f °C\n", sg.SecondaryOutletTemp())
	fmt.Printf("\tHeat Transfer Rate: %.2f MW\n", sg.HeatTransferRate())
	fmt.Printf("\tSteam Flow Rate: %.2f kg/s\n", sg.SteamFlowRate())
}

func (sg *SteamGenerator) Update(s *simworks.Simulator) (map[string]interface{}, error) {

	if sg.primaryLoop == nil || sg.secondaryLoop == nil {
		return nil, errors.New("steam generator must be connected to primary and secondary loops")
	}

	sg.BaseComponent.Update(s)
	// Update primary inlet temperature based on reactor core heat
	sg.primaryInletTemp = sg.primaryLoop.HotLegTemperature() // expect 325˚C
	sg.primaryOutletTemp = sg.primaryInletTemp - (35.0)      // aiming for 290˚C if hot leg is 325˚C

	sg.secondaryOutletTemp = 285.0                                     // aiming for 285˚C, but depends on primary inlet temp
	sg.secondaryInletTemp = sg.secondaryLoop.FeedwaterTemperatureOut() // 80˚C if feedheaters are on or around 40˚C if not

	sg.heatTransferRate = 0.0 // TODO: calculate this using primary flow volume and temp differential between primary and secondary
	sg.steamFlowRate = 0.0    // TODO: calculate this using heat transfer rate (and knowledge of reality)

	return sg.Status(), nil
}
