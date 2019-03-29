package main

import (
	"fmt"
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/logging"
	"github.com/prashantgupta24/activity-tracker/pkg/tracker"
)

func main() {

	logger := logging.New()

	logger.Infof("starting activity tracker")

	frequency := 12 //value always in seconds

	activityTracker := &tracker.Instance{
		Frequency: frequency,
		LogLevel:  logging.Info,
	}

	//This starts the tracker for all services
	heartbeatCh := activityTracker.Start()

	//if you only want to track certain services, you can use StartWithServices
	//heartbeatCh := activityTracker.StartWithServices(service.MouseClickHandler(), service.MouseCursorHandler())

	timeToKill := time.NewTicker(time.Second * 60)

	for {
		select {
		case heartbeat := <-heartbeatCh:
			if !heartbeat.WasAnyActivity {
				logger.Infof("no activity detected in the last %v seconds\n\n\n", int(frequency))
			} else {
				logger.Info(heartbeat.ActivityMap)
				logger.Infof("activity detected in the last %v seconds.", int(frequency))
				logger.Infof("Activity type:\n")
				for activity, time := range heartbeat.ActivityMap {
					//logger.Infof("%v ---> %v\n", activity.ActivityType, time)
					logger.Infof("activity : %v times: %v\n", activity, time)
				}
				fmt.Printf("\n\n\n")
			}
		case <-timeToKill.C:
			logger.Infof("time to kill app")
			activityTracker.Quit()
			return
		}
	}
}
