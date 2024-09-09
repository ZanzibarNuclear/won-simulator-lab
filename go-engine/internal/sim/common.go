package sim

import (
	"math/rand"
)

type Environment struct {
	Weather string
}

func NewEnvironment() *Environment {
	return &Environment{
		Weather: "sunny",
	}
}

func (e *Environment) String() string {
	return e.Weather
}

const ROOM_TEMPERATURE = 20
const TURBINE_MAX_RPM = 3600

// common run durations, assuming 1-minute iterations
const HOUR_OF_MINUTES = 60
const DAY_OF_MINUTES = HOUR_OF_MINUTES * 24
const WEEK_OF_MINUTES = DAY_OF_MINUTES * 7
const YEAR_OF_MINUTES = WEEK_OF_MINUTES * 52

func generateRandomID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
