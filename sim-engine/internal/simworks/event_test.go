package simworks

import (
	"testing"
	"time"
)

func TestNewImmediateEvent(t *testing.T) {
	e := NewImmediateEvent("test")
	if e.Code != "test" {
		t.Errorf("Expected code to be 'test', got %s", e.Code)
	}
	if e.Status != "pending" {
		t.Errorf("Expected status to be 'pending', got %s", e.Status)
	}
	if e.Immediate != true {
		t.Errorf("Expected immediate to be true, got %t", e.Immediate)
	}
}

func TestNewAdjustmentEvent(t *testing.T) {
	e := NewAdjustmentEvent("test", 10.0)
	if e.Code != "test" {
		t.Errorf("Expected code to be 'test', got %s", e.Code)
	}
	if e.Status != "pending" {
		t.Errorf("Expected status to be 'pending', got %s", e.Status)
	}
	if e.TargetValue != 10.0 {
		t.Errorf("Expected target value to be 10.0, got %f", e.TargetValue)
	}
	if e.Immediate {
		t.Errorf("Expected immediate to be false, got true")
	}
}

func TestScheduleAt(t *testing.T) {
	e := NewImmediateEvent("test")
	startTime := time.Now().Add(time.Second * 10)
	e.ScheduleAt(startTime)

	if e.StartMoment.IsZero() {
		t.Errorf("Expected start moment to be set, got %v", e.StartMoment)
	}
	if e.StartMoment != startTime {
		t.Errorf("Expected start moment to be %v, got %v", startTime, e.StartMoment)
	}
}

func TestForDuration(t *testing.T) {
	e := NewAdjustmentEvent("test", 10.0)
	duration := time.Second * 10
	e.ForDuration(duration)

	if e.EndMoment().IsZero() {
		t.Errorf("Expected end moment to be set, got %v", e.EndMoment())
	}
	if e.EndMoment().Sub(e.StartMoment) != duration {
		t.Errorf("Expected duration to be %v, got %v", duration, e.EndMoment().Sub(e.StartMoment))
	}
}

func TestAdjustmentEventScheduledWithDuration(t *testing.T) {
	// Create a new adjustment event
	e := NewAdjustmentEvent("test_adjustment", 100.0)

	// Get the current time
	now := time.Now()

	// Schedule the event for an hour from now
	startTime := now.Add(time.Hour)
	duration := time.Minute
	e.ScheduleAt(startTime).ForDuration(duration)

	// Check if the start time is set correctly
	if !e.StartMoment.Equal(startTime) {
		t.Errorf("Expected start time to be %v, got %v", startTime, e.StartMoment)
	}

	// Check if the end time is set correctly (start time + duration)
	expectedEndTime := startTime.Add(duration)
	if !e.EndMoment().Equal(expectedEndTime) {
		t.Errorf("Expected end time to be %v, got %v", expectedEndTime, e.EndMoment())
	}

	// Check if the duration is correct
	if e.EndMoment().Sub(e.StartMoment) != duration {
		t.Errorf("Expected duration to be %v, got %v", duration, e.EndMoment().Sub(e.StartMoment))
	}

	startTime = now.Add(time.Second * 10)
	duration = time.Hour
	e.ForDuration(duration).ScheduleAt(startTime)

	if !e.StartMoment.Equal(startTime) {
		t.Errorf("Expected start time to be %v, got %v", startTime, e.StartMoment)
	}
	if e.EndMoment().Sub(e.StartMoment) != duration {
		t.Errorf("Expected duration to be %v, got %v", duration, e.EndMoment().Sub(e.StartMoment))
	}
}

func TestSettingDurationOnImmediateEvent(t *testing.T) {
	e := NewImmediateEvent("test")
	e.ForDuration(time.Second * 10)

	if !e.EndMoment().IsZero() {
		t.Errorf("Expected end moment not to be set, got %v", e.EndMoment())
	}
}

func TestAdjustmentEvent_ZeroDuration(t *testing.T) {
	e := NewAdjustmentEvent("test", 100.0)
	e.ForDuration(0)

	if !e.EndMoment().IsZero() {
		t.Errorf("Zero duration should set EndMoment to StartMoment, got %v", e.EndMoment())
	}
}

func TestAdjustmentEvent_NegativeDuration(t *testing.T) {
	e := NewAdjustmentEvent("test", 100.0)
	e.ForDuration(time.Duration(-1))

	if e.EndMoment().IsZero() {
		t.Errorf("Negative duration to be ignored, got %v", e.EndMoment())
	}
}

func TestEventStatusChanges(t *testing.T) {
	e := NewImmediateEvent("test_status")

	// Test SetComplete
	e.SetComplete()
	if e.Status != "completed" {
		t.Errorf("Expected status to be 'completed', got %s", e.Status)
	}

	// Test SetInProgress
	e.SetInProgress()
	if e.Status != "in_progress" {
		t.Errorf("Expected status to be 'in_progress', got %s", e.Status)
	}

	// Test SetPending
	e.SetPending()
	if e.Status != "pending" {
		t.Errorf("Expected status to be 'pending', got %s", e.Status)
	}

	// Test SetCanceled
	e.SetCanceled()
	if e.Status != "canceled" {
		t.Errorf("Expected status to be 'canceled', got %s", e.Status)
	}
}
