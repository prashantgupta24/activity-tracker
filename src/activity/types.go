package activity

import "time"

type ActivityTracker struct {
	TimeToCheck time.Duration
}

type Heartbeat struct {
	IsActivity bool
	Time       time.Time
}
