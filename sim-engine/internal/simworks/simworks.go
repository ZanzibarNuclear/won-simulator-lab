package simworks

import (
	"fmt"
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
	Events         []*Event
	InactiveEvents []*Event
	EventHandler   EventHandler
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
		Events:         []*Event{},
		InactiveEvents: []*Event{},
	}
}

// SetEventHandler sets the event handler for the simulator
func (s *Simulator) SetEventHandler(handler EventHandler) {
	s.EventHandler = handler
}

func (s *Simulator) Run(seconds int) {
	results := make(map[string]interface{})
	for i := 0; i < seconds; i++ {
		s.Clock.Tick()

		s.ProgressDueEvents()

		for _, component := range s.Components {
			status, _ := component.Update(s)
			results[component.ID()] = status
		}

		if s.EventHandler != nil {
			s.ReviewPendingEvents()
		}
		s.TidyUpEvents()
	}

	results["general"] = "TODO"
	// TODO: decide what to do with results - that's a lot of data
}

func (s *Simulator) Step() {
	s.Run(1)
}

func (s *Simulator) ProgressDueEvents() {
	for _, event := range s.Events {
		if event.IsPending() && event.IsDue(s.CurrentMoment()) {
			event.SetInProgress()
		}
	}
}

func (s *Simulator) ReviewPendingEvents() {
	for i, event := range s.Events {
		if event.IsPending() && event.IsDue(s.CurrentMoment()) {
			s.EventHandler.ProcessEvent(event)
			fmt.Printf("Event status post processing: %s\n", s.Events[i].Status)
		}
	}
}

func (s *Simulator) TidyUpEvents() {
	n := 0
	for _, event := range s.Events {
		if event.IsComplete() || event.IsCanceled() {
			s.InactiveEvents = append(s.InactiveEvents, event)
		} else {
			s.Events[n] = event
			n++
		}
	}
	s.Events = s.Events[:n]
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

func (s *Simulator) QueueEvent(e *Event) {
	s.Events = append(s.Events, e)
}

func (s *Simulator) PrintStatus() {
	fmt.Printf("=== Simulator %s is running. ===\n  %s\n\n", s.Name, s.Purpose)
	for _, component := range s.Components {
		component.Print()
	}
	for _, event := range s.Events {
		fmt.Println(event.State())
	}
}
