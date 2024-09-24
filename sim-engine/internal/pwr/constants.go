package pwr

// physical constants and typical operational settings
const (
	ROOM_TEMPERATURE = 20.0 // in Celsius

	PUMP_ON_PRESSURE         = 1.0   // MPa
	PUMP_ON_FLOW_RATE        = 20.0  // in m³/s
	PUMP_ON_HEAT             = 100.0 // in kW
	PUMP_OFF_PRESSURE        = 0.0   // MPa
	PUMP_OFF_FLOW_RATE       = 0.0
	PUMP_OFF_HEAT            = 0.0
	MAX_BORON_RATE_OF_CHANGE = 0.083  // ppm/second
	MAX_BORON_CONCENTRATION  = 2500.0 // ppm
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

var Config = map[string]map[string]float64{
	"primary_loop": PrimaryLoopConfig,
	"common": {
		"room_temperature": ROOM_TEMPERATURE,
	},
}

// event codes
const (
	Event_pl_pumpSwitch         = "primary_loop.cooling_pump.switch"
	Event_pl_boronConcentration = "primary_loop.cvcs.boron_concentration_target"
)
