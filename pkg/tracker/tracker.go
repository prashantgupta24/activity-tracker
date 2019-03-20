package tracker

import (
	"log"
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/service"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

const (
	preHeartbeatTime = time.Millisecond * 100
)

func (tracker *Instance) Start() (heartbeatCh chan *Heartbeat, quit chan struct{}) {

	//register service handlers
	tracker.registerHandlers(&service.MouseClickHandler{}, &service.MouseCursorHandler{},
		&service.ScreenChangeHandler{})

	//returned channels
	heartbeatCh = make(chan *Heartbeat, 1)
	quit = make(chan struct{})

	go func(tracker *Instance) {
		timeToCheck := tracker.TimeToCheck
		//tickers
		tickerHeartbeat := time.NewTicker(time.Second * timeToCheck)
		tickerWorker := time.NewTicker(time.Second*timeToCheck - preHeartbeatTime)

		activities := makeActivityMap()

		for {
			select {
			case <-tickerWorker.C:
				log.Printf("tracker worker working at %v\n", time.Now())
				//time to ping all registered service handlers
				//doing it the non-blocking sender way
				for _, serviceHandler := range tracker.serviceHandlers {
					select {
					case serviceHandler <- struct{}{}:
					default:
						//service is blocked
					}
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
				//close all service handlers for a clean exit
				for _, serviceHandler := range tracker.serviceHandlers {
					close(serviceHandler)
				}
				return
			}
		}
	}(tracker)

	return heartbeatCh, quit
}

func makeActivityMap() map[*activity.Type]time.Time {
	activityMap := make(map[*activity.Type]time.Time)
	return activityMap
}

func (tracker *Instance) registerHandlers(services ...service.Instance) {

	if len(tracker.serviceHandlers) == 0 { //checking for multiple registration attempts
		tracker.activityCh = make(chan *activity.Type, len(services)) // number based on types of activities being tracked

		for _, service := range services {
			tickerCh := service.Start(tracker.activityCh)
			tracker.serviceHandlers = append(tracker.serviceHandlers, tickerCh)
		}
	}
}
