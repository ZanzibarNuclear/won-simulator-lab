package pwr

// SteamTableEntry represents a single entry in the steam table
type SteamTableEntry struct {
	Temperature    float64 // in °C
	Pressure       float64 // in MPa
	SpecificVolume float64 // in m³/kg
	Enthalpy       float64 // in kJ/kg
}

// SteamTable is a slice of SteamTableEntry representing the steam table
var SteamTable = []SteamTableEntry{
	{Temperature: 0.01, Pressure: 0.000611, SpecificVolume: 0.001000, Enthalpy: 0.00},
	{Temperature: 20, Pressure: 0.002339, SpecificVolume: 0.001002, Enthalpy: 83.96},
	{Temperature: 50, Pressure: 0.012349, SpecificVolume: 0.001012, Enthalpy: 209.33},
	{Temperature: 100, Pressure: 0.101325, SpecificVolume: 0.001043, Enthalpy: 419.04},
	{Temperature: 120, Pressure: 0.198540, SpecificVolume: 0.001060, Enthalpy: 503.71},
	{Temperature: 140, Pressure: 0.361390, SpecificVolume: 0.001079, Enthalpy: 589.13},
	{Temperature: 160, Pressure: 0.617800, SpecificVolume: 0.001100, Enthalpy: 675.47},
	{Temperature: 180, Pressure: 1.002600, SpecificVolume: 0.001123, Enthalpy: 762.81},
	{Temperature: 200, Pressure: 1.554900, SpecificVolume: 0.001148, Enthalpy: 851.24},
	{Temperature: 220, Pressure: 2.319600, SpecificVolume: 0.001175, Enthalpy: 940.87},
	{Temperature: 240, Pressure: 3.344000, SpecificVolume: 0.001204, Enthalpy: 1031.80},
	{Temperature: 260, Pressure: 4.688000, SpecificVolume: 0.001236, Enthalpy: 1124.20},
	{Temperature: 280, Pressure: 6.412000, SpecificVolume: 0.001272, Enthalpy: 1218.30},
	{Temperature: 300, Pressure: 8.581000, SpecificVolume: 0.001312, Enthalpy: 1314.30},
	{Temperature: 320, Pressure: 11.270000, SpecificVolume: 0.001357, Enthalpy: 1412.60},
	{Temperature: 340, Pressure: 14.586000, SpecificVolume: 0.001410, Enthalpy: 1513.80},
	{Temperature: 360, Pressure: 18.651000, SpecificVolume: 0.001475, Enthalpy: 1619.00},
	{Temperature: 374.14, Pressure: 22.064000, SpecificVolume: 0.003155, Enthalpy: 2099.30},
}

// GetSteamProperties returns the steam properties for a given temperature
func GetSteamProperties(temperature float64) (SteamTableEntry, bool) {
	for _, entry := range SteamTable {
		if entry.Temperature == temperature {
			return entry, true
		}
	}
	return SteamTableEntry{}, false
}

// InterpolateSteamProperties returns interpolated steam properties for a given temperature
func InterpolateSteamProperties(temperature float64) SteamTableEntry {
	if temperature <= SteamTable[0].Temperature {
		return SteamTable[0]
	}
	if temperature >= SteamTable[len(SteamTable)-1].Temperature {
		return SteamTable[len(SteamTable)-1]
	}

	var lowIndex int
	for i, entry := range SteamTable {
		if entry.Temperature > temperature {
			lowIndex = i - 1
			break
		}
	}

	low := SteamTable[lowIndex]
	high := SteamTable[lowIndex+1]
	ratio := (temperature - low.Temperature) / (high.Temperature - low.Temperature)

	return SteamTableEntry{
		Temperature:    temperature,
		Pressure:       low.Pressure + ratio*(high.Pressure-low.Pressure),
		SpecificVolume: low.SpecificVolume + ratio*(high.SpecificVolume-low.SpecificVolume),
		Enthalpy:       low.Enthalpy + ratio*(high.Enthalpy-low.Enthalpy),
	}
}

// InterpolateFromGivenPressure returns interpolated steam properties for a given pressure
func InterpolateFromGivenPressure(pressure float64) SteamTableEntry {
	if pressure <= SteamTable[0].Pressure {
		return SteamTable[0]
	}
	if pressure >= SteamTable[len(SteamTable)-1].Pressure {
		return SteamTable[len(SteamTable)-1]
	}

	var lowIndex int
	for i, entry := range SteamTable {
		if entry.Pressure > pressure {
			lowIndex = i - 1
			break
		}
	}

	low := SteamTable[lowIndex]
	high := SteamTable[lowIndex+1]
	ratio := (pressure - low.Pressure) / (high.Pressure - low.Pressure)

	return SteamTableEntry{
		Temperature:    low.Temperature + ratio*(high.Temperature-low.Temperature),
		Pressure:       pressure,
		SpecificVolume: low.SpecificVolume + ratio*(high.SpecificVolume-low.SpecificVolume),
		Enthalpy:       low.Enthalpy + ratio*(high.Enthalpy-low.Enthalpy),
	}
}
