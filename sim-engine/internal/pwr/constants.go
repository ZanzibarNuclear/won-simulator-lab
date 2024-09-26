package pwr

// physical constants and typical component settings

var PrimaryLoopConfig = map[string]float64{
	"pump_on_pressure":         1.0,
	"pump_on_flow_rate":        20.0,
	"pump_on_heat":             100.0,
	"pump_off_pressure":        0.0,
	"pump_off_flow_rate":       0.0,
	"pump_off_heat":            0.0,
	"max_boron_rate_of_change": 0.083,
	"max_boron_concentration":  2500.0,
}

var PressurizerConfig = map[string]float64{
	"target_pressure":                 15.5,
	"target_temperature":              345.0,
	"heater_high_power":               1500.0,
	"heater_low_power":                50.0,
	"spray_flow_rate":                 10.0,
	"relief_valve_threshold_pressure": 18.0,
	"relief_valve_drop_rate":          0.5, // MPa/s
}

var Config = map[string]map[string]float64{
	"common": {
		"room_temperature":     20.0,
		"atmospheric_pressure": 0.101325,
	},
	"primary_loop": PrimaryLoopConfig,
	"pressurizer":  PressurizerConfig,
}
