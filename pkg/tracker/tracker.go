package tracker

import (
	"log"
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/service"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

const (
	preHeartbeatTime = time.Millisecond * 10
)

func (tracker *Instance) Start() (heartbeatCh chan *Heartbeat, quit chan struct{}) {

	//register service handlers
	tracker.registerHandlers(service.MouseClickHandler, service.MouseCursorHandler)

	//returned channels
	heartbeatCh = make(chan *Heartbeat, 1)
	quit = make(chan struct{})

	go func(tracker *Instance, heartbeatCh chan *Heartbeat, quit chan struct{}) {

		timeToCheck := tracker.TimeToCheck
		//tickers
		tickerHeartbeat := time.NewTicker(time.Second * timeToCheck)
		tickerWorker := time.NewTicker(time.Second*timeToCheck - preHeartbeatTime)

		activities := makeActivityMap()

		for {
			select {
			case <-tickerWorker.C:
				log.Printf("tracker worker working at %v\n", time.Now())
				for i := 0; i < len(tracker.services); i++ {
					tracker.workerTickerCh <- struct{}{}
				}
			case <-tickerHeartbeat.C:
				log.Printf("tracker heartbeat checking at %v\n", time.Now())
				var heartbeat *Heartbeat
				if len(activities) == 0 {
					//log.Printf("no activity detected in the last %v seconds ...\n", int(timeToCheck))
					heartbeat = &Heartbeat{
						IsActivity: false,
						Activity:   nil,
						Time:       time.Now(),
					}
				} else {
					//log.Printf("activity detected in the last %v seconds ...\n", int(timeToCheck))
					heartbeat = &Heartbeat{
						IsActivity: true,
						Activity:   activities,
						Time:       time.Now(),
					}

				}
				heartbeatCh <- heartbeat
				activities = makeActivityMap() //reset the activities map
			case activity := <-tracker.activityCh:
				activities[activity] = time.Now()
				//log.Printf("activity received: %#v\n", activity)
			case <-quit:
				log.Printf("stopping activity tracker\n")
				for _, quitHandler := range tracker.services {
					quitHandler <- struct{}{}
				}
				//robotgo.StopEvent()
				return
			}
		}
	}(tracker, heartbeatCh, quit)

	return heartbeatCh, quit
}

func makeActivityMap() map[*activity.Type]time.Time {
	activityMap := make(map[*activity.Type]time.Time)
	return activityMap
}

func (tracker *Instance) registerHandlers(services ...func(tickerCh chan struct{},
	clickComm chan *activity.Type) (quit chan struct{})) {

	tracker.activityCh = make(chan *activity.Type, len(services)) // number based on types of activities being tracked
	tracker.workerTickerCh = make(chan struct{}, len(services))   //this is for all the services, instead of each having their own

	for _, service := range services {
		quitCh := service(tracker.workerTickerCh, tracker.activityCh)
		tracker.services = append(tracker.services, quitCh)
	}
}
