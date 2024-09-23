package simworks

import (
	"testing"
	"time"
)

func SpawnClock() *SimClock {
	return SpawnClockAt(testStart)
}

func SpawnClockAt(simStart time.Time) *SimClock {
	return NewClock(simStart)
}

var testStart = time.Date(1971, time.March, 1, 18, 20, 37, 0, time.UTC)

func TestClockSeconds(t *testing.T) {
	testClock := SpawnClock()

	// Tick the clock for 13 months
	for i := 0; i < 59; i++ {
		testClock.Tick()
	}

	if testClock.Ticks() != 59 {
		t.Errorf("Expected Ticks to be 59, but got %d", testClock.Ticks())
	}

	if testClock.SimNow().Sub(testStart) != time.Duration(59)*time.Second {
		t.Errorf("Expected SimNow to be %v, but got %v", testStart.Add(time.Duration(59)*time.Second), testClock.SimNow())
	}

	if testClock.FormatNow() != "01-Mar-1971 18:21:36" {
		t.Errorf("Expected FormatNow to be %v, but got %v", "01-Mar-1971 18:21:36", testClock.FormatNow())
	}
}

func TestAdvancingClockMinutes(t *testing.T) {
	testClock := SpawnClock()

	// Advance the clock by 50 minutes (3000 ticks)
	for i := 0; i < 3000; i++ {
		testClock.Tick()
	}

	if testClock.Ticks() != 3000 {
		t.Errorf("Expected Ticks to be 3000, but got %d", testClock.Ticks())
	}

	expectedTime := testStart.Add(time.Duration(3000) * time.Second)
	if testClock.SimNow() != expectedTime {
		t.Errorf("Expected SimNow to be %v, but got %v", expectedTime, testClock.SimNow())
	}

	if testClock.FormatNow() != "01-Mar-1971 19:10:37" {
		t.Errorf("Expected FormatNow to be %v, but got %v", "01-Mar-1971 19:10:37", testClock.FormatNow())
	}
}

func TestAdvancingClockHours(t *testing.T) {
	testClock := SpawnClock()

	// Advance the clock by 7 hours (25200 ticks)
	for i := 0; i < 25200; i++ {
		testClock.Tick()
	}

	if testClock.Ticks() != 25200 {
		t.Errorf("Expected Ticks to be 3660, but got %d", testClock.Ticks())
	}

	expectedTime := testStart.Add(time.Duration(25200) * time.Second)
	if testClock.SimNow() != expectedTime {
		t.Errorf("Expected SimNow to be %v, but got %v", expectedTime, testClock.SimNow())
	}

	if testClock.FormatNow() != "02-Mar-1971 01:20:37" {
		t.Errorf("Expected FormatNow to be %v, but got %v", "02-Mar-1971 01:20:37", testClock.FormatNow())
	}
}

func TestAdvancingClockDays(t *testing.T) {
	testClock := SpawnClock()

	// Advance the clock by 32 days (2764800 ticks)
	for i := 0; i < 2764800; i++ {
		testClock.Tick()
	}

	if testClock.Ticks() != 2764800 {
		t.Errorf("Expected Ticks to be 86460, but got %d", testClock.Ticks())
	}

	expectedTime := testStart.Add(time.Duration(2764800) * time.Second)
	if testClock.SimNow() != expectedTime {
		t.Errorf("Expected SimNow to be %v, but got %v", expectedTime, testClock.SimNow())
	}

	if testClock.FormatNow() != "02-Apr-1971 18:20:37" {
		t.Errorf("Expected FormatNow to be %v, but got %v", "02-Apr-1971 18:20:37", testClock.FormatNow())
	}
}

func TestAdvancingClockYears(t *testing.T) {
	testClock := SpawnClock()

	// Advance clock by 1 year
	for i := 0; i < 31536000; i++ {
		testClock.Tick()
	}

	if testClock.Ticks() != 31536000 {
		t.Errorf("Expected Ticks to be 31536000, but got %d", testClock.Ticks())
	}

	expectedTime := testStart.Add(time.Duration(365*24*60*60) * time.Second)
	if testClock.SimNow() != expectedTime {
		t.Errorf("Expected SimNow to be %v, but got %v", expectedTime, testClock.SimNow())
	}

	if testClock.FormatNow() != "29-Feb-1972 18:20:37" {
		t.Errorf("1972 was a leap year; Expected %v, but got %v", "29-Feb-1972 18:20:37", testClock.FormatNow())
	}
}
