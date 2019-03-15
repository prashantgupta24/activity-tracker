package activity

import "time"

type ActivityType string

const (
	MOUSE_CURSOR_MOVEMENT ActivityType = "cursor-move"
	MOUSE_LEFT_CLICK      ActivityType = "left-mouse-click"
)

type Activity struct {
	ActivityType ActivityType
}

type ActivityTracker struct {
	TimeToCheck time.Duration
}

type Heartbeat struct {
	IsActivity bool
	Activity   *Activity
	Time       time.Time
}
