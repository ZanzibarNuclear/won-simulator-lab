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

var SecondaryLoopConfig = map[string]float64{
	"ideal_steam_temperature":       285.0,
	"mssv_pressure_threshold":       8.0,
	"base_feedwater_temperature":    40.0,
	"heated_feedwater_temperature":  80.0,
	"feedheater_step_up":            1.25,
	"feedheater_step_down":          2.5,
	"feedwater_flow_rate_target":    20.0,
	"feedwater_flow_rate_step_up":   0.8,
	"feedwater_flow_rate_step_down": 1.5,
}

var SteamTurbineConfig = map[string]float64{
	"max_rpm":        3600,
	"efficiency":     0.85,
	"blade_diameter": 0.6,
}

var Config = map[string]map[string]float64{
	"common": {
		"room_temperature":     20.0,
		"atmospheric_pressure": 0.101325,
		"gravity":              9.81,
	},
	"primary_loop":   PrimaryLoopConfig,
	"secondary_loop": SecondaryLoopConfig,
	"pressurizer":    PressurizerConfig,
	"steam_turbine":  SteamTurbineConfig,
}
