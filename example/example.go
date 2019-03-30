package main

import (
	"fmt"
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/logging"
	"github.com/prashantgupta24/activity-tracker/pkg/tracker"
)

func main() {

	logger := logging.New()

	frequency := 60 //value always in seconds

	activityTracker := &tracker.Instance{
		Frequency: frequency,
		LogLevel:  logging.Debug,
	}

	//This starts the tracker for all handlers
	heartbeatCh := activityTracker.Start()

	//if you only want to track certain handlers, you can use StartWithhandlers
	//heartbeatCh := activityTracker.StartWithHanders(handler.MouseClickHandler(), handler.MouseCursorHandler())

	timeToKill := time.NewTicker(time.Second * 60)

	logger.Infof("starting activity tracker with %v second frequency ...", frequency)

	for {
		select {
		case heartbeat := <-heartbeatCh:
			if !heartbeat.WasAnyActivity {
				logger.Infof("no activity detected in the last %v seconds\n\n\n", int(frequency))
			} else {
				logger.Infof("activity detected in the last %v seconds.", int(frequency))
				logger.Infof("Activity type:\n")
				for activityType, times := range heartbeat.ActivityMap {
					logger.Infof("activityType : %v times: %v\n", activityType, len(times))
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
