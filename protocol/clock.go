package protocol

import (
	"time"
)

type Clock interface {
	Now() time.Time
	NewElectionTimer(time.Duration) *time.Timer
	StopElectionTimer(*time.Timer)
}

type physicalClock struct{}

var PhysicalClock = physicalClock{}

func (physicalClock) Now() time.Time {
	return time.Now()
}

func (physicalClock) NewElectionTimer(t time.Duration) *time.Timer {
	return time.NewTimer(t)
}

func (physicalClock) StopElectionTimer(t *time.Timer) {
	t.Stop()
}
