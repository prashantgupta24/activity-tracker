package tracker

import (
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/service"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

//Instance is an instance of the tracker
type Instance struct {
	Frequency  int
	LogLevel   string
	LogFormat  string
	activityCh chan *activity.Type
	quit       chan struct{}
	services   map[activity.Type]service.Instance
}

//Heartbeat is the data packet sent from the tracker to the user
type Heartbeat struct {
	WasAnyActivity bool
	Activity       map[*activity.Type]time.Time //activity type with its time
	Time           time.Time                    //heartbeat time
}
