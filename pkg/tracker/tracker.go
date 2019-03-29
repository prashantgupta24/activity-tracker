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

//StartWithServices starts the tracker with a set of services
func (tracker *Instance) StartWithServices(services ...service.Instance) (heartbeatCh chan *Heartbeat) {
	logger := logging.NewLoggerLevelFormat(tracker.LogLevel, tracker.LogFormat)

	//register service handlers
	tracker.registerHandlers(logger, services...)

	//returned channel
	heartbeatCh = make(chan *Heartbeat, 1)

	tracker.quit = make(chan struct{})

	go func(logger *log.Logger, tracker *Instance) {
		trackerLog := logger.WithFields(log.Fields{
			"method": "activity-tracker",
		})

		timeToCheck := time.Duration(tracker.Frequency)
		//tickers
		tickerHeartbeat := time.NewTicker(timeToCheck * time.Second)
		var tickerWorker *time.Ticker
		if timeToCheck >= 10 {
			tickerWorker = time.NewTicker((timeToCheck / 5 * time.Second) - preHeartbeatTime)
		} else {
			tickerWorker = time.NewTicker(timeToCheck*time.Second - preHeartbeatTime)
		}

		activityMap := makeActivityMap()

		for {
			select {
			case <-tickerWorker.C:
				trackerLog.Infof("tracker worker working")
				//time to trigger all registered services
				for _, service := range tracker.services {
					service.Trigger()
				}
			case <-tickerHeartbeat.C:
				trackerLog.Debugln("tracker heartbeat checking")
				var heartbeat *Heartbeat
				if len(activityMap) == 0 {
					logger.Debugf("no activity detected in the last %v seconds ...\n", int(timeToCheck))
					heartbeat = &Heartbeat{
						WasAnyActivity: false,
						ActivityMap:    nil,
						Time:           time.Now(),
					}
				} else {
					trackerLog.Debugf("activity detected in the last %v seconds ...\n", int(timeToCheck))
					heartbeat = &Heartbeat{
						WasAnyActivity: true,
						ActivityMap:    activityMap,
						Time:           time.Now(),
					}
				}
				heartbeatCh <- heartbeat
				activityMap = makeActivityMap() //reset the activityMap map
				trackerLog.Debugln("**************** END OF CHECK ********************")
			case activity := <-tracker.activityCh:
				timeNow := time.Now()
				activityMap[activity.Type] = append(activityMap[activity.Type], timeNow)
				trackerLog.Debugf("activity received: \n%#v\n", activity)
			case <-tracker.quit:
				trackerLog.Infof("stopping activity tracker\n")
				//close all services for a clean exit
				for _, service := range tracker.services {
					service.Close()
				}
				close(heartbeatCh)
				return
			}
		}
	}(logger, tracker)

	return heartbeatCh
}

//Quit the tracker app
func (tracker *Instance) Quit() {
	tracker.quit <- struct{}{}
}

//Start the tracker with all possible services
func (tracker *Instance) Start() (heartbeatCh chan *Heartbeat) {
	return tracker.StartWithServices(getAllServiceHandlers()...)
}

func getAllServiceHandlers() []service.Instance {
	return []service.Instance{
		service.MouseClickHandler(), service.MouseCursorHandler(),
		service.ScreenChangeHandler(),
	}
}

func makeActivityMap() (activityMap map[activity.Type][]time.Time) {
	activityMap = make(map[activity.Type][]time.Time)
	return activityMap
}

func (tracker *Instance) registerHandlers(logger *log.Logger, services ...service.Instance) {

	tracker.services = make(map[activity.Type]service.Instance)
	tracker.activityCh = make(chan *activity.Instance, len(services)) // number based on types of activities being tracked

	for _, service := range services {
		if _, ok := tracker.services[service.Type()]; !ok { //duplicate registration prevention
			tracker.services[service.Type()] = service
			service.Start(logger, tracker.activityCh)
		}
	}
}
