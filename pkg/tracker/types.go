package tracker

import (
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/service"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

type Instance struct {
	Frequency  int
	activityCh chan *activity.Type
	quit       chan struct{}
	services   map[service.Instance]bool
}

type Heartbeat struct {
	IsActivity bool
	Activity   map[*activity.Type]time.Time
	Time       time.Time
}
