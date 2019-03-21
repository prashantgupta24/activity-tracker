package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/logging"
	"github.com/prashantgupta24/activity-tracker/pkg/tracker"
)

func main() {

	logger := logging.New()

	logger.Infof("starting activity tracker")

	frequency := 5 //value always in seconds

	activityTracker := &tracker.Instance{
		Frequency: frequency,
	}

	//This starts the tracker for all services
	heartbeatCh := activityTracker.Start()

	//if you only want to track certain services, you can use StartWithServices
	//heartbeatCh := activityTracker.StartWithServices(service.MouseClickHandler(), service.MouseCursorHandler())

	timeToKill := time.NewTicker(time.Second * 30)

	for {
		select {
		case heartbeat := <-heartbeatCh:
			if !heartbeat.IsActivity {
				log.Infof("no activity detected in the last %v seconds\n\n", int(frequency))
			} else {
				log.Infof("activity detected in the last %v seconds. ", int(frequency))
				log.Infof("Activity type:\n")
				for activity, time := range heartbeat.Activity {
					log.Infof("%v ---> %v\n", activity.ActivityType, time)
				}
				log.Println()
			}
		case <-timeToKill.C:
			log.Infof("time to kill app")
			activityTracker.Quit()
			return
		}
	}
}
