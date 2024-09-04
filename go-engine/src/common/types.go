package common

type Environment struct {
	Weather string
}

const ROOM_TEMPERATURE = 20
const TURBINE_MAX_RPM = 3600

// common run durations, assuming 1-minute iterations
const HOUR_OF_MINUTES = 60
const DAY_OF_MINUTES = HOUR_OF_MINUTES * 24
const WEEK_OF_MINUTES = DAY_OF_MINUTES * 7
const YEAR_OF_MINUTES = WEEK_OF_MINUTES * 52