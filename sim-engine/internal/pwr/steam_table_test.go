package pwr

import (
	"testing"

	"worldofnuclear.com/internal/simworks"
)

func TestInterpolateSteamProperties(t *testing.T) {
	testCases := []struct {
		temperature float64
		expected    SteamTableEntry
		tolerance   float64
	}{
		{
			temperature: 130,
			expected: SteamTableEntry{
				Temperature:    130,
				Pressure:       0.279965,
				SpecificVolume: 0.0010695,
				Enthalpy:       546.42,
			},
			tolerance: 0.0001,
		},
		{
			temperature: 290,
			expected: SteamTableEntry{
				Temperature:    290,
				Pressure:       7.49649,
				SpecificVolume: 0.001292,
				Enthalpy:       1266.3,
			},
			tolerance: 0.0001,
		},
	}

	for _, tc := range testCases {
		result := InterpolateSteamProperties(tc.temperature)

		if !simworks.AlmostEqual(result.Temperature, tc.expected.Temperature, tc.tolerance) ||
			!simworks.AlmostEqual(result.Pressure, tc.expected.Pressure, tc.tolerance) ||
			!simworks.AlmostEqual(result.SpecificVolume, tc.expected.SpecificVolume, tc.tolerance) ||
			!simworks.AlmostEqual(result.Enthalpy, tc.expected.Enthalpy, tc.tolerance) {
			t.Errorf("For temperature %.2f°C, expected %+v, but got %+v", tc.temperature, tc.expected, result)
		}
	}
}

func TestInterpolateFromGivenPressure(t *testing.T) {
	testCases := []struct {
		pressure  float64
		expected  SteamTableEntry
		tolerance float64
	}{
		{
			pressure: 2.0,
			expected: SteamTableEntry{
				Temperature:    211.64116,
				Pressure:       2.0,
				SpecificVolume: 0.001163,
				Enthalpy:       903.40988,
			},
			tolerance: 0.001,
		},
		{
			pressure: 7.5,
			expected: SteamTableEntry{
				Temperature:    290.032272,
				Pressure:       7.5,
				SpecificVolume: 0.001175,
				Enthalpy:       1266.4549,
			},
			tolerance: 0.001,
		},
	}

	for _, tc := range testCases {
		result := InterpolateFromGivenPressure(tc.pressure)

		if !simworks.AlmostEqual(result.Temperature, tc.expected.Temperature, tc.tolerance) ||
			!simworks.AlmostEqual(result.Pressure, tc.expected.Pressure, tc.tolerance) ||
			!simworks.AlmostEqual(result.SpecificVolume, tc.expected.SpecificVolume, tc.tolerance) ||
			!simworks.AlmostEqual(result.Enthalpy, tc.expected.Enthalpy, tc.tolerance) {
			t.Errorf("For pressure %.4f MPa, expected %+v, but got %+v", tc.pressure, tc.expected, result)
		}
	}
}
