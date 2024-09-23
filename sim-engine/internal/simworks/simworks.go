package simworks

import (
	"time"
)

// Simulator represents a simulation system
type Simulator struct {
	ID           string
	Name         string
	Purpose      string
	CreationDate time.Time
	Clock        *SimClock
	Environment  *Environment
}

// NewSimulator creates a new simulator
func NewSimulator(id, name, purpose string) *Simulator {
	return &Simulator{
		ID:           id,
		Name:         name,
		Purpose:      purpose,
		CreationDate: time.Now(),
		Clock:        NewClock(time.Now()),
		Environment:  NewEnvironment(),
	}
}

func (s *Simulator) Run(seconds int) {
	for i := 0; i < seconds; i++ {
		s.Clock.Tick()
	}
}

func (s *Simulator) RunForABit(days, hours, minutes, seconds int) {
	duration := days*ONE_DAY_IN_SECONDS + hours*ONE_HOUR_IN_SECONDS + minutes*ONE_MINUTE_IN_SECONDS + seconds
	s.Run(duration)
}
