package pwr

// physical constants and typical operational settings
const (
	ROOM_TEMPERATURE     = 20.0     // in Celsius
	ATMOSPHERIC_PRESSURE = 0.101325 // MPa

	PUMP_ON_PRESSURE         = 1.0   // MPa
	PUMP_ON_FLOW_RATE        = 20.0  // in m³/s
	PUMP_ON_HEAT             = 100.0 // in kW
	PUMP_OFF_PRESSURE        = 0.0   // MPa
	PUMP_OFF_FLOW_RATE       = 0.0
	PUMP_OFF_HEAT            = 0.0
	MAX_BORON_RATE_OF_CHANGE = 0.083  // ppm/second
	MAX_BORON_CONCENTRATION  = 2500.0 // ppm

	TARGET_PRESSURE                 = 15.5   // MPa
	TARGET_TEMPERATURE              = 345.0  // °C
	HEATER_HIGH_POWER               = 1500.0 // kW
	HEATER_LOW_POWER                = 50.0   // kW, enough to hold steady
	SPRAY_FLOW_RATE                 = 10.0   // kg/s
	RELIEF_VALVE_FLOW               = 50.0   // kg/s
	RELIEF_VALVE_THRESHOLD_PRESSURE = 17.0   // MPa
	RELIEF_VALVE_DROP_RATE          = 0.5    // MPa/s TODO: switch to drop rate
)

var PrimaryLoopConfig = map[string]float64{
	"pump_on_pressure":         PUMP_ON_PRESSURE,
	"pump_on_flow_rate":        PUMP_ON_FLOW_RATE,
	"pump_on_heat":             PUMP_ON_HEAT,
	"pump_off_pressure":        PUMP_OFF_PRESSURE,
	"pump_off_flow_rate":       PUMP_OFF_FLOW_RATE,
	"pump_off_heat":            PUMP_OFF_HEAT,
	"max_boron_rate_of_change": MAX_BORON_RATE_OF_CHANGE,
	"max_boron_concentration":  MAX_BORON_CONCENTRATION,
}

var PressurizerConfig = map[string]float64{
	"target_pressure":                 TARGET_PRESSURE,
	"target_temperature":              TARGET_TEMPERATURE,
	"heater_high_power":               HEATER_HIGH_POWER,
	"heater_low_power":                HEATER_LOW_POWER,
	"spray_flow_rate":                 SPRAY_FLOW_RATE,
	"relief_valve_flow":               RELIEF_VALVE_FLOW,
	"relief_valve_threshold_pressure": RELIEF_VALVE_THRESHOLD_PRESSURE,
	"relief_valve_drop_rate":          RELIEF_VALVE_DROP_RATE,
}

var Config = map[string]map[string]float64{
	"common": {
		"room_temperature":     ROOM_TEMPERATURE,
		"atmospheric_pressure": ATMOSPHERIC_PRESSURE,
	},
	"primary_loop": PrimaryLoopConfig,
	"pressurizer":  PressurizerConfig,
}
