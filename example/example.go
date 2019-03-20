package main

import (
	"fmt"
	"time"

	"github.com/prashantgupta24/activity-tracker/pkg/tracker"
)

func main() {
	fmt.Println("starting activity tracker")

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
				fmt.Printf("no activity detected in the last %v seconds\n\n", int(frequency))
			} else {
				fmt.Printf("activity detected in the last %v seconds. ", int(frequency))
				fmt.Printf("Activity type:\n")
				for activity, time := range heartbeat.Activity {
					fmt.Printf("%v ---> %v\n", activity.ActivityType, time)
				}
				fmt.Println()
			}
		case <-timeToKill.C:
			fmt.Println("time to kill app")
			activityTracker.Quit()
			return
		}
	}
}
