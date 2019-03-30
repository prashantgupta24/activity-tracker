package main

import (
	"fmt"
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/logging"
	"github.com/prashantgupta24/activity-tracker/pkg/tracker"
)

func main() {

	logger := logging.New()

	heartbeatFrequency := 60 //value always in seconds
	workerFrequency := 5     //seconds

	activityTracker := &tracker.Instance{
		HeartbeatFrequency: heartbeatFrequency,
		WorkerFrequency:    workerFrequency,
		LogLevel:           logging.Info,
	}

	//This starts the tracker for all handlers
	heartbeatCh := activityTracker.Start()

	//if you only want to track certain handlers, you can use StartWithhandlers
	//heartbeatCh := activityTracker.StartWithHanders(handler.MouseClickHandler(), handler.MouseCursorHandler())

	timeToKill := time.NewTicker(time.Second * 120)

	logger.Infof("starting activity tracker with %vs heartbeat and %vs worker frequency...", heartbeatFrequency, workerFrequency)

	for {
		select {
		case heartbeat := <-heartbeatCh:
			if !heartbeat.WasAnyActivity {
				logger.Infof("no activity detected in the last %v seconds\n\n\n", int(heartbeatFrequency))
			} else {
				logger.Infof("activity detected in the last %v seconds.", int(heartbeatFrequency))
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
