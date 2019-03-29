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

<<<<<<< Updated upstream
//Heartbeat is the data packet sent from the tracker to the user
=======
/*Heartbeat is the data packet sent from the tracker to the user.

WasAnyActivity tells if there was any activity within that time frame
If there was, then the ActivityMap will tell you what type of activity
it was and at what times it occured.

The Time field is the time of the Heartbeat sent (not to be confused with
the activity time, which is the time the activity occured within the time frame)
*/
>>>>>>> Stashed changes
type Heartbeat struct {
	WasAnyActivity bool
	ActivityMap    map[*activity.Type][]time.Time //activity type with its times
	Time           time.Time                      //heartbeat time
}
