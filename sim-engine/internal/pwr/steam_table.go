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
	{Temperature: 120, Pressure: 0.200000, SpecificVolume: 0.001061, Enthalpy: 504.71},
	{Temperature: 140, Pressure: 0.316228, SpecificVolume: 0.001081, Enthalpy: 593.84},
	{Temperature: 160, Pressure: 0.476101, SpecificVolume: 0.001090, Enthalpy: 632.20},
	{Temperature: 180, Pressure: 0.681292, SpecificVolume: 0.001101, Enthalpy: 672.79},
	{Temperature: 200, Pressure: 0.933378, SpecificVolume: 0.001114, Enthalpy: 715.01},
	{Temperature: 220, Pressure: 1.234000, SpecificVolume: 0.001129, Enthalpy: 758.35},
	{Temperature: 240, Pressure: 1.584893, SpecificVolume: 0.001146, Enthalpy: 802.46},
	{Temperature: 260, Pressure: 1.988190, SpecificVolume: 0.001165, Enthalpy: 847.02},
	{Temperature: 280, Pressure: 2.445200, SpecificVolume: 0.001186, Enthalpy: 891.76},
	{Temperature: 300, Pressure: 2.957350, SpecificVolume: 0.001209, Enthalpy: 936.40},
	{Temperature: 320, Pressure: 3.525800, SpecificVolume: 0.001234, Enthalpy: 980.76},
	{Temperature: 250, Pressure: 3.973920, SpecificVolume: 0.001252, Enthalpy: 1085.75},
	{Temperature: 300, Pressure: 8.587840, SpecificVolume: 0.001395, Enthalpy: 1345.11},
	{Temperature: 350, Pressure: 16.529160, SpecificVolume: 0.001642, Enthalpy: 1648.1},
	{Temperature: 374.14, Pressure: 22.064000, SpecificVolume: 0.003155, Enthalpy: 2099.3},
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
