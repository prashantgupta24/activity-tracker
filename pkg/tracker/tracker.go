package tracker

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/logging"
	"github.com/prashantgupta24/activity-tracker/internal/pkg/service"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

const (
	preHeartbeatTime = time.Millisecond * 100
)

func (tracker *Instance) StartWithServices(services ...service.Instance) (heartbeatCh chan *Heartbeat) {
	logger := logging.NewLoggerLevelFormat(tracker.LogLevel, tracker.LogFormat)

	//register service handlers
	tracker.registerHandlers(logger, services...)

	//returned channels
	heartbeatCh = make(chan *Heartbeat, 1)
	tracker.quit = make(chan struct{})

	go func(tracker *Instance) {
		timeToCheck := time.Duration(tracker.Frequency)
		//tickers
		tickerHeartbeat := time.NewTicker(time.Second * timeToCheck)
		tickerWorker := time.NewTicker(time.Second*timeToCheck - preHeartbeatTime)

		activities := makeActivityMap()

		for {
			select {
			case <-tickerWorker.C:
				logger.Debugln("tracker worker working")
				//time to trigger all registered services
				for service := range tracker.services {
					service.Trigger()
				}
			case <-tickerHeartbeat.C:
				logger.Debugln("tracker heartbeat checking")
				var heartbeat *Heartbeat
				if len(activities) == 0 {
					logger.Debugf("no activity detected in the last %v seconds ...\n", int(timeToCheck))
					heartbeat = &Heartbeat{
						IsActivity: false,
						Activity:   nil,
						Time:       time.Now(),
					}
				} else {
					logger.Debugf("activity detected in the last %v seconds ...\n", int(timeToCheck))
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
				logger.Debugf("activity received: %#v\n", activity)
			case <-tracker.quit:
				logger.Infof("stopping activity tracker\n")
				//close all services for a clean exit
				for service := range tracker.services {
					service.Close()
				}
				return
			}
		}
	}(tracker)

	return heartbeatCh
}

func (tracker *Instance) Quit() {
	tracker.quit <- struct{}{}
}

func (tracker *Instance) Start() (heartbeatCh chan *Heartbeat) {
	return tracker.StartWithServices(service.MouseClickHandler(), service.MouseCursorHandler(),
		service.ScreenChangeHandler())
}

func makeActivityMap() map[*activity.Type]time.Time {
	activityMap := make(map[*activity.Type]time.Time)
	return activityMap
}

func (tracker *Instance) registerHandlers(logger *log.Logger, services ...service.Instance) {

	tracker.services = make(map[service.Instance]bool)
	tracker.activityCh = make(chan *activity.Type, len(services)) // number based on types of activities being tracked

	for _, service := range services {
		service.Start(logger, tracker.activityCh)
		if _, ok := tracker.services[service]; !ok { //duplicate registration prevention
			tracker.services[service] = true
		}

	}
}
