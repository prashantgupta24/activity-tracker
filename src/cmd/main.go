package main

import (
	"fmt"
	"time"

	"github.com/activitytracker/src/activity"
)

func main() {
	fmt.Println("starting activity tracker")

	activityCh, quitActivityDetector := activity.StartTracker(5)

	timeToKill := time.NewTicker(time.Second * 30)

	for {
		select {
		case isActivity := <-activityCh:
			if !isActivity {
				fmt.Println("no activity detected")
			} else {
				fmt.Println("activity detected")
			}
		case <-timeToKill.C:
			fmt.Println("time to kill app")
			quitActivityDetector <- struct{}{}
			return
		}
	}
}
