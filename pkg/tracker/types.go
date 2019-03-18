package tracker

import (
	"time"

	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

type Instance struct {
	TimeToCheck    time.Duration
	activityCh     chan *activity.Type
	workerTickerCh chan struct{}
	services       []chan struct{}
}

type Heartbeat struct {
	IsActivity bool
	Activity   map[*activity.Type]time.Time
	Time       time.Time
}
