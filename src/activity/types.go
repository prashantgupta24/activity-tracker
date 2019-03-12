package activity

import "time"

type ActivityTracker struct {
	TimeToCheck time.Duration
	Debug       bool
}

type Heartbeat struct {
	isActivity bool
	time       time.Time
}
