package simworks

import "time"

type Event struct {
	Code        string        // key to identify the event
	Status      string        // "pending", "in_progress", "completed", or "canceled"
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
func NewImmediateEvent(code string) *Event {
	return &Event{
		Code:      code,
		Status:    "pending",
		Immediate: true,
	}
}

func NewImmediateEventBool(code string, value bool) *Event {
	var targetValue float64
	if value {
		targetValue = 1.0
	} else {
		targetValue = 0.0
	}
	return &Event{
		Code:        code,
		Status:      "pending",
		Immediate:   true,
		TargetValue: targetValue,
	}
}

// NewAdjustmentEvent creates a new adjustment event with the given code and target value
func NewAdjustmentEvent(code string, targetValue float64) *Event {
	return &Event{
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

// Status returns a map containing the current status of the event
func (e *Event) State() map[string]interface{} {
	status := map[string]interface{}{
		"code":         e.Code,
		"status":       e.Status,
		"immediate":    e.Immediate,
		"target_value": e.TargetValue,
	}

	if !e.StartMoment.IsZero() {
		status["start_moment"] = e.StartMoment
	}

	if e.Duration != 0 {
		status["duration"] = e.Duration
	}

	return status
}

func (e *Event) SetComplete() {
	e.Status = "completed"
}

func (e *Event) IsComplete() bool {
	return e.Status == "completed"
}

func (e *Event) SetInProgress() {
	e.Status = "in_progress"
}

func (e *Event) IsInProgress() bool {
	return e.Status == "in_progress"
}

func (e *Event) SetPending() {
	e.Status = "pending"
}

func (e *Event) IsPending() bool {
	return e.Status == "pending"
}

func (e *Event) SetCanceled() {
	e.Status = "canceled"
}

func (e *Event) IsCanceled() bool {
	return e.Status == "canceled"
}

func (e *Event) IsScheduled() bool {
	return !e.StartMoment.IsZero()
}

func (e *Event) IsDue(moment time.Time) bool {
	if !e.IsScheduled() {
		return true
	} else if moment.Equal(e.StartMoment) || moment.After(e.StartMoment) {
		return true
	}
	return false
}

func (e *Event) Truthy() bool {
	return e.TargetValue != 0
}

type EventHandler interface {
	ProcessEvent(e *Event)
}
