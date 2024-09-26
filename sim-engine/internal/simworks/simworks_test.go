package simworks

import (
	"fmt"
	"testing"
)

func TestSimWorks(t *testing.T) {
	s := NewSimulator("Tess the Tester", "Test Simulator")

	if len(s.ID) != 8 {
		t.Errorf("Expected ID to be 8 characters, got %s", s.ID)
	}

	if s.Name != "Tess the Tester" {
		t.Errorf("Expected Name to be 'Tess the Tester', got %s", s.Name)
	}

	if s.Purpose != "Test Simulator" {
		t.Errorf("Expected Purpose to be 'Test Simulator', got %s", s.Purpose)
	}

	if s.CreationDate.IsZero() {
		t.Errorf("Expected CreationDate to be non-zero, got %s", s.CreationDate)
	}

	if s.Clock == nil {
		t.Errorf("Expected Clock to be non-nil, got %v", s.Clock)
	}

	if s.Environment == nil {
		t.Errorf("Expected Environment to be non-nil, got %v", s.Environment)
	}

	if len(s.Components) != 0 {
		t.Errorf("Expected Components to be empty, got %v", s.Components)
	}

	if len(s.ComponentIndex) != 0 {
		t.Errorf("Expected ComponentIndex to be empty, got %v", s.ComponentIndex)
	}
}

const (
	Event_ignore    = "ignore-this-event"
	Event_count     = "count-this-event"
	Event_times_ten = "times-ten-this-event"
)

type EventHandlerThingy struct {
	count int
}

func (e *EventHandlerThingy) ProcessEvent(event *Event) {
	fmt.Printf("Processing event: %s\n", event.Code)
	switch event.Code {
	case Event_count:
		e.count++
		event.SetComplete()
	case Event_times_ten:
		e.count *= 10
		event.SetComplete()
	}
	fmt.Printf("Event status is: %s\n", event.Status)
}

func TestSimulator_EventHandling(t *testing.T) {
	s := NewSimulator("Tess the Tester", "Test Simulator")

	handler := &EventHandlerThingy{}
	s.SetEventHandler(handler)

	s.QueueEvent(NewImmediateEvent(Event_count))
	s.QueueEvent(NewImmediateEvent(Event_ignore))
	s.QueueEvent(NewImmediateEvent(Event_count))
	s.QueueEvent(NewImmediateEvent(Event_times_ten))

	s.Step()

	if handler.count != 20 {
		t.Errorf("Expected count to be 20, got %d", handler.count)
	}

	if len(s.Events) != 1 {
		t.Errorf("Expected one unprocessed event, got %v", s.Events)
	}

	if len(s.InactiveEvents) != 3 {
		t.Errorf("Expected 3 inactive events, got %d", len(s.InactiveEvents))
	}

	fmt.Printf("Events: %v\n", s.Events)
	fmt.Printf("Inactive Events: %v\n", s.InactiveEvents)
}
