package simworks

import (
	"time"
)

// SimClock represents a clock that keep track of the simulation time
// It holds a "start" time that represents the point from which the
// simulation begins running, and the number of ticks that have
// occurred since then. A tick might represent a second of simulation
// time, but it's ambiguous for flexibility.

type SimClock struct {
	Start time.Time
	ticks int
}

const (
	ONE_MINUTE_IN_SECONDS    = 60
	ONE_HOUR_IN_SECONDS      = 60 * ONE_MINUTE_IN_SECONDS
	ONE_DAY_IN_SECONDS       = 24 * ONE_HOUR_IN_SECONDS
	ONE_WEEK_IN_SECONDS      = 7 * ONE_DAY_IN_SECONDS
	ONE_YEAR_IN_SECONDS      = 52*ONE_WEEK_IN_SECONDS + ONE_DAY_IN_SECONDS
	ONE_LEAP_YEAR_IN_SECONDS = 52*ONE_WEEK_IN_SECONDS + 2*ONE_DAY_IN_SECONDS
)

// NewClock creates a new clock
func NewClock(simStart time.Time) *SimClock {
	return &SimClock{
		Start: simStart.Truncate(time.Second),
	}
}

func NewDefaultClock() *SimClock {
	return NewClock(time.Now())
}

// The simulation clock only goes forward, one tick at a time
func (c *SimClock) Tick() {
	c.ticks++
}

// Ticks returns the number of ticks
func (c *SimClock) Ticks() int {
	return c.ticks
}

// Elapsed returns the elapsed time
func (c *SimClock) SimNow() time.Time {
	return c.Start.Add(time.Duration(c.ticks) * time.Second)
}

// Elapsed returns the elapsed time
func (c *SimClock) FormatNow() string {
	return c.SimNow().Format("02-Jan-2006 15:04:05")
}
