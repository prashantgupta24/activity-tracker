package main

import (
	"fmt"
	"time"

	"github.com/prashantgupta24/activity-tracker/src/activity"
)

func main() {
	fmt.Println("starting activity tracker")

	timeToCheck := 5

	activityTracker := &activity.ActivityTracker{
		TimeToCheck: time.Duration(timeToCheck),
	}
	heartbeatCh, quitActivityTracker := activityTracker.Start()

	timeToKill := time.NewTicker(time.Second * 30)

	for {
		select {
		case heartbeat := <-heartbeatCh:
			if !heartbeat.IsActivity {
				fmt.Printf("no activity detected in the last %v seconds\n", int(timeToCheck))
			} else {
				fmt.Printf("activity detected in the last %v seconds. ", int(timeToCheck))
				fmt.Printf("Activity type %#v\n", heartbeat.Activity)
			}
		case <-timeToKill.C:
			fmt.Println("time to kill app")
			quitActivityTracker <- struct{}{}
			return
		}
	}
}
