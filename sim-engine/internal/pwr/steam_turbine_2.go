package pwr

import (
	"fmt"
	"math"
	"strings"
)

/*
Advanced model of the steam turbine. Might come in useful when considering whether to
level up models for each component.

This model takes into account the following factors:
- Moisture content
- Feedwater heaters
- Condenser type
- Turbine stages
- Extraction points
- Efficiency
*/

// TurbineStage represents a single stage in the steam turbine
type TurbineStage struct {
	Name            string
	InletPressure   float64 // MPa
	OutletPressure  float64 // MPa
	InletTemp       float64 // °C
	OutletTemp      float64 // °C
	MoistureContent float64 // fraction
}

// TurbineParams holds input parameters for turbine calculations
type TurbineParams struct {
	ThrottlePressure         float64
	ThrottleTemp             float64
	GeneratorOutput          float64
	MoistureSeparationStages int
	FeedwaterHeaters         int
	CondenserType            string
}

// TurbineResults holds the calculated results
type TurbineResults struct {
	Stages              []TurbineStage
	ExtractionPressures []float64 // MPa
	TotalEfficiency     float64
	SteamFlow           float64 // kg/s
	CondenserPressure   float64 // MPa
	GeneratorOutput     float64 // MW
}

// Constants for calculations
const (
	maxMoistureContent   = 0.12
	typicalMSRTempRise   = 40.0
	waterCooledCondenser = "water-cooled"
	airCooledCondenser   = "air-cooled"
)

// pressureRatios holds typical pressure ratios for PWR turbine stages
var pressureRatios = map[string]float64{
	"HP": 0.25,
	"IP": 0.20,
	"LP": 0.004,
}

func calculatePWRTurbine(params TurbineParams) TurbineResults {
	var results TurbineResults

	// Convert MPa to Pa for internal calculations
	throttlePressurePa := params.ThrottlePressure * 1e6

	// Calculate saturation temperature if not provided
	saturationTemp := math.Min(params.ThrottleTemp,
		234.0+44.0*math.Log10(params.ThrottlePressure))

	// Initialize stage calculations
	var stages []TurbineStage
	currentPressure := throttlePressurePa
	currentTemp := saturationTemp

	// HP Stage
	hpOutletPressure := currentPressure * pressureRatios["HP"]
	hpOutletTemp := 180.0   // Typical for PWR HP outlet
	moistureContent := 0.12 // Typical moisture at HP outlet

	stages = append(stages, TurbineStage{
		Name:            "HP",
		InletPressure:   currentPressure / 1e6,
		OutletPressure:  hpOutletPressure / 1e6,
		InletTemp:       currentTemp,
		OutletTemp:      hpOutletTemp,
		MoistureContent: moistureContent,
	})

	// Moisture Separator Reheater stages
	for i := 0; i < params.MoistureSeparationStages; i++ {
		pressureDrop := 0.95
		if i > 0 {
			pressureDrop = 0.90
		}
		inletPressure := hpOutletPressure * pressureDrop
		outletTemp := currentTemp + typicalMSRTempRise
		outletPressure := inletPressure * 0.95

		stages = append(stages, TurbineStage{
			Name:            fmt.Sprintf("MSR%d", i+1),
			InletPressure:   inletPressure / 1e6,
			OutletPressure:  outletPressure / 1e6,
			InletTemp:       hpOutletTemp,
			OutletTemp:      outletTemp,
			MoistureContent: 0.001,
		})
		currentPressure = outletPressure
		currentTemp = outletTemp
	}

	// LP Stages
	lpStages := 4 // Typical number of LP stages
	for i := 0; i < lpStages; i++ {
		stagePressureRatio := math.Pow(pressureRatios["LP"], 1.0/float64(lpStages))
		outletPressure := currentPressure * stagePressureRatio
		outletTemp := 60 + (currentTemp-60)*math.Pow(stagePressureRatio, 0.3)
		moistureContent = math.Min(maxMoistureContent, moistureContent+0.03)

		stages = append(stages, TurbineStage{
			Name:            fmt.Sprintf("LP%d", i+1),
			InletPressure:   currentPressure / 1e6,
			OutletPressure:  outletPressure / 1e6,
			InletTemp:       currentTemp,
			OutletTemp:      outletTemp,
			MoistureContent: moistureContent,
		})
		currentPressure = outletPressure
		currentTemp = outletTemp
	}

	// Calculate condenser pressure
	var condenserPressure float64
	if params.CondenserType == waterCooledCondenser {
		condenserPressure = 5000 // Pa
	} else {
		condenserPressure = 8000 // Pa
	}

	// Calculate extraction points
	var extractionPressures []float64
	totalPressureRatio := condenserPressure / throttlePressurePa

	for i := 0; i < params.FeedwaterHeaters; i++ {
		ratio := math.Pow(totalPressureRatio, float64(i+1)/float64(params.FeedwaterHeaters+1))
		extractionPressures = append(extractionPressures, throttlePressurePa*ratio/1e6) // Convert to MPa
	}

	// Calculate efficiency
	baseEfficiency := 0.75
	msrBonus := 0.02 * float64(params.MoistureSeparationStages)
	heaterBonus := 0.002 * float64(params.FeedwaterHeaters)
	totalEfficiency := math.Min(0.86, baseEfficiency+msrBonus+heaterBonus)

	// Calculate steam flow rate
	enthalpyDrop := 800.0 // kJ/kg
	steamFlow := (params.GeneratorOutput * 1000) / (enthalpyDrop * totalEfficiency)

	// Compile results
	results = TurbineResults{
		Stages:              stages,
		ExtractionPressures: extractionPressures,
		TotalEfficiency:     totalEfficiency,
		SteamFlow:           steamFlow,
		CondenserPressure:   condenserPressure / 1000,
		GeneratorOutput:     params.GeneratorOutput,
	}

	return results
}

func printResults(results TurbineResults) {
	fmt.Println("\nPWR Steam Turbine Analysis")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Total Efficiency: %.1f%%\n", results.TotalEfficiency*100)
	fmt.Printf("Steam Flow Rate: %.1f kg/s\n", results.SteamFlow)
	fmt.Printf("Generator Output: %.0f MW\n", results.GeneratorOutput)
	fmt.Printf("Condenser Pressure: %.2f kPa\n", results.CondenserPressure)

	fmt.Println("\nTurbine Stages:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-6s %-15s %-15s %-15s %-15s %-10s\n",
		"Stage", "In Press(MPa)", "Out Press(MPa)", "In Temp(°C)", "Out Temp(°C)", "Moisture %")
	fmt.Println(strings.Repeat("-", 80))

	for _, stage := range results.Stages {
		fmt.Printf("%-6s %15.3f %15.3f %15.1f %15.1f %10.1f\n",
			stage.Name, stage.InletPressure, stage.OutletPressure,
			stage.InletTemp, stage.OutletTemp, stage.MoistureContent*100)
	}

	fmt.Println("\nFeedwater Heater Extraction Points:")
	fmt.Println(strings.Repeat("-", 50))
	for i, pressure := range results.ExtractionPressures {
		fmt.Printf("Extraction %d: %.3f MPa\n", i+1, pressure)
	}
}

// func main() {
// 	fmt.Println("PWR Nuclear Power Plant Steam Turbine Calculator")
// 	fmt.Println(strings.Repeat("=", 50))

// 	// Get input parameters with typical PWR values as defaults
// 	params := TurbineParams{
// 		ThrottlePressureMPa:      readFloat("Enter main steam pressure (MPa)", 7.0),
// 		ThrottleTempC:            readFloat("Enter main steam temperature (°C)", 286.0),
// 		GeneratorOutputMW:        readFloat("Enter required generator output (MW)", 1100.0),
// 		MoistureSeparationStages: readInt("Enter number of moisture separator/reheater stages", 1),
// 		FeedwaterHeaters:         readInt("Enter number of feedwater heaters", 7),
// 		CondenserType:            waterCooledCondenser,
// 	}

// 	fmt.Println("\nCooling options:")
// 	fmt.Println("1. water-cooled")
// 	fmt.Println("2. air-cooled")
// 	coolingChoice := readInt("Select cooling type (1-2)", 1)
// 	if coolingChoice == 2 {
// 		params.CondenserType = airCooledCondenser
// 	}

// 	// Calculate and display results
// 	results := calculatePWRTurbine(params)
// 	printResults(results)
// }
