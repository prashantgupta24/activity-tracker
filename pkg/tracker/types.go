package tracker

import (
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/handler"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
	"github.com/prashantgupta24/activity-tracker/pkg/system"
)

//Instance is an instance of the tracker
type Instance struct {
	HeartbeatInterval int //the interval at which you want the heartbeat (in seconds, default 60s)
	WorkerInterval    int //the interval at which you want the checks to happen within a heartbeat (in seconds, default 5s)
	LogLevel          string
	LogFormat         string
	state             *system.State //maintains the state of the system
	isTest            bool          //only for testing purposes
	activityCh        chan *activity.Instance
	quit              chan struct{}
	handlers          map[activity.Type]handler.Instance
}

/*Heartbeat is the data packet sent from the tracker to the user.

WasAnyActivity tells if there was any activity within that time frame
If there was, then the ActivityMap will tell you what type of activity
it was and at what times it occured.

The Time field is the time of the Heartbeat sent (not to be confused with
the activity time, which is the time the activity occured within the time frame)
*/
type Heartbeat struct {
	WasAnyActivity bool
	ActivityMap    map[activity.Type][]time.Time //activity type with its times
	Time           time.Time                     //heartbeat time
}
