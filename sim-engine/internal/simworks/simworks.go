package simworks

import (
	"time"
)

// Simulator represents a simulation system
type Simulator struct {
	ID             string
	Name           string
	Purpose        string
	CreationDate   time.Time
	Clock          *SimClock
	Environment    *Environment
	Components     []SimComponent
	ComponentIndex map[string]SimComponent
	Events         []Event
}

// NewSimulator creates a new simulator
func NewSimulator(name, purpose string) *Simulator {
	return &Simulator{
		ID:             GenerateRandomID(8),
		Name:           name,
		Purpose:        purpose,
		CreationDate:   time.Now(),
		Clock:          NewClock(time.Now()),
		Environment:    NewEnvironment(),
		Components:     []SimComponent{},
		ComponentIndex: make(map[string]SimComponent),
		Events:         []Event{},
	}
}

func (s *Simulator) Run(seconds int) {
	for i := 0; i < seconds; i++ {
		s.Clock.Tick()
		for _, component := range s.Components {
			component.Update(s)
		}
	}
}

func (s *Simulator) Step() {
	s.Run(1)
}

func (s *Simulator) CurrentMoment() time.Time {
	return s.Clock.SimNow()
}

func (s *Simulator) RunForABit(days, hours, minutes, seconds int) {
	duration := days*ONE_DAY_IN_SECONDS + hours*ONE_HOUR_IN_SECONDS + minutes*ONE_MINUTE_IN_SECONDS + seconds
	s.Run(duration)
}

func (s *Simulator) AddComponent(c SimComponent) {
	s.Components = append(s.Components, c)
	s.ComponentIndex[c.ID()] = c
}

func (s *Simulator) QueueEvent(e Event) {
	s.Events = append(s.Events, e)
}
