package main

import (
	"fmt"
	"time"

	"github.com/prashantgupta24/activity-tracker/src/activity"
)

func main() {
	fmt.Println("starting activity tracker")

	activityTracker := &activity.ActivityTracker{
		TimeToCheck: 5,
	}
	heartbeatCh, quitActivityTracker := activityTracker.Start()

	timeToKill := time.NewTicker(time.Second * 30)

	for {
		select {
		case heartbeat := <-heartbeatCh:
			if !heartbeat.IsActivity {
				fmt.Println("no activity detected")
			} else {
				fmt.Println("activity detected")
			}
		case <-timeToKill.C:
			fmt.Println("time to kill app")
			quitActivityTracker <- struct{}{}
			return
		}
	}
}
