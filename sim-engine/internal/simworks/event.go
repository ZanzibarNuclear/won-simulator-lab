package simworks

import "time"

type Event struct {
	Code        string        // key to identify the event
	Status      string        // "pending", "in_progress", or "completed"
	Immediate   bool          // for changes that happen in a single tick
	StartMoment time.Time     // to schedule a future event
	Duration    time.Duration // positive for events that have a defined duration; negative not allowed
	TargetValue float64       // only non-zero for events that have a defined target value
	FromValue   float64       // can be useful for interpolation, or sense of direction of change
}

func (e *Event) EndMoment() time.Time {
	return e.StartMoment.Add(e.Duration)
}

// NewImmediateEvent creates a new immediate event with the given code
func NewImmediateEvent(code string) Event {
	return Event{
		Code:      code,
		Status:    "pending",
		Immediate: true,
	}
}

// NewAdjustmentEvent creates a new adjustment event with the given code and target value
func NewAdjustmentEvent(code string, targetValue float64) Event {
	return Event{
		Code:        code,
		Status:      "pending",
		TargetValue: targetValue,
	}
}

// Have event start at a specific time
func (e *Event) ScheduleAt(t time.Time) *Event {
	e.StartMoment = t
	return e
}

// Have event run for a given duration
func (e *Event) ForDuration(duration time.Duration) *Event {
	if e.Immediate {
		return e
	}
	e.Duration = duration
	return e
}

func (e *Event) SetComplete() {
	e.Status = "completed"
}

func (e *Event) SetInProgress() {
	e.Status = "in_progress"
}

func (e *Event) SetPending() {
	e.Status = "pending"
}

func (e *Event) SetCanceled() {
	e.Status = "canceled"
}
