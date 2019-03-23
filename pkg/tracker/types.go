package tracker

import (
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/service"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

type Instance struct {
	Frequency  int
	LogLevel   string
	LogFormat  string
	activityCh chan *activity.Type
	quit       chan struct{}
	services   map[activity.Type]service.Instance
}

type Heartbeat struct {
	IsActivity bool
	Activity   map[*activity.Type]time.Time //activity type with its time
	Time       time.Time                    //heartbeat time
}
