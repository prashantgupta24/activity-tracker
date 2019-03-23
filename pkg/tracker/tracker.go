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

	go func(logger *log.Logger, tracker *Instance) {
		trackerLog := logger.WithFields(log.Fields{
			"method": "activity-tracker",
		})
		timeToCheck := time.Duration(tracker.Frequency)
		//tickers
		tickerHeartbeat := time.NewTicker(time.Second * timeToCheck)
		tickerWorker := time.NewTicker(time.Second*timeToCheck - preHeartbeatTime)

		activities := makeActivityMap()

		for {
			select {
			case <-tickerWorker.C:
				trackerLog.Debugln("tracker worker working")
				//time to trigger all registered services
				for _, service := range tracker.services {
					service.Trigger()
				}
			case <-tickerHeartbeat.C:
				trackerLog.Debugln("tracker heartbeat checking")
				var heartbeat *Heartbeat
				if len(activities) == 0 {
					logger.Debugf("no activity detected in the last %v seconds ...\n", int(timeToCheck))
					heartbeat = &Heartbeat{
						IsActivity: false,
						Activity:   nil,
						Time:       time.Now(),
					}
				} else {
					trackerLog.Debugf("activity detected in the last %v seconds ...\n", int(timeToCheck))
					heartbeat = &Heartbeat{
						IsActivity: true,
						Activity:   activities,
						Time:       time.Now(),
					}

				}
				heartbeatCh <- heartbeat
				activities = makeActivityMap() //reset the activities map
				trackerLog.Debugln("**************** END OF CHECK ********************")
			case activity := <-tracker.activityCh:
				activities[activity] = time.Now()
				trackerLog.Debugf("activity received: \n%#v\n", activity)
			case <-tracker.quit:
				trackerLog.Infof("stopping activity tracker\n")
				//close all services for a clean exit
				for _, service := range tracker.services {
					service.Close()
				}
				return
			}
		}
	}(logger, tracker)

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

	tracker.services = make(map[activity.Type]service.Instance)
	tracker.activityCh = make(chan *activity.Type, len(services)) // number based on types of activities being tracked

	for _, service := range services {
		service.Start(logger, tracker.activityCh)
		if _, ok := tracker.services[service.Type()]; !ok { //duplicate registration prevention
			tracker.services[service.Type()] = service
		}
	}
}
