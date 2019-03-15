package activity

import "time"

type ActivityType string

const (
	MOUSE_CURSOR_MOVEMENT ActivityType = "cursor-move"
	MOUSE_LEFT_CLICK      ActivityType = "left-mouse-click"
	SCREEN_CHANGE         ActivityType = "screen-change"
)

type Activity struct {
	ActivityType ActivityType
}

type ActivityTracker struct {
	TimeToCheck    time.Duration
	activityCh     chan *Activity
	workerTickerCh chan struct{}
	services       []chan struct{}
}

type Heartbeat struct {
	IsActivity bool
	Activity   map[*Activity]time.Time
	Time       time.Time
}
