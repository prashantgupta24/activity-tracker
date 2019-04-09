package tracker

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/handler"
	"github.com/prashantgupta24/activity-tracker/internal/pkg/logging"
	"github.com/prashantgupta24/activity-tracker/pkg/system"
)

const (
	preHeartbeatTime = time.Millisecond * 100
	//heartbeat (seconds)
	minHInterval     = 60
	maxHInterval     = 300
	defaultHInterval = 60

	//worker (seconds)
	minWInterval     = 4
	maxWInterval     = minHInterval
	defaultWInterval = minHInterval
)

//StartWithHandlers starts the tracker with a set of handlers
func (tracker *Instance) StartWithHandlers(handlers ...handler.Instance) (heartbeatCh chan *Heartbeat) {
	logger := logging.NewLoggerLevelFormat(tracker.LogLevel, tracker.LogFormat)

	//register handlers
	tracker.registerHandlers(logger, handlers...)

	//returned channel
	heartbeatCh = make(chan *Heartbeat, 1)

	tracker.quit = make(chan struct{})
	tracker.state = &system.State{
		IsSystemSleep: false,
	}

	go func(logger *log.Logger, tracker *Instance) {
		trackerLog := logger.WithFields(log.Fields{
			"method": "activity-tracker",
		})

		//instantiating ticker intervals
		heartbeatInterval, workerInterval := tracker.validateIntervals()

		tickerHeartbeat := time.NewTicker(heartbeatInterval * time.Second)
		tickerWorker := time.NewTicker(workerInterval*time.Second - preHeartbeatTime)

		activityMap := makeActivityMap()

		for {
			select {
			case <-tickerWorker.C:
				trackerLog.Debugln("tracker worker working")
				//time to trigger all registered handlers
				for _, handler := range tracker.handlers {
					handler.Trigger(*tracker.state) //sending a copy of the state
				}
			case <-tickerHeartbeat.C:
				trackerLog.Debugln("tracker heartbeat checking")
				var heartbeat *Heartbeat
				if len(activityMap) == 0 {
					logger.Debugf("no activity detected in the last %v seconds ...\n", int(heartbeatInterval))
					heartbeat = &Heartbeat{
						WasAnyActivity: false,
						ActivityMap:    nil,
						Time:           time.Now(),
					}
				} else {
					trackerLog.Debugf("activity detected in the last %v seconds ...\n", int(heartbeatInterval))
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
				if activity.State != nil {
					trackerLog.Debugf("received state change request for %v activity", activity.Type)
					tracker.state = activity.State
				}
				timeNow := time.Now()
				activityMap[activity.Type] = append(activityMap[activity.Type], timeNow)
				trackerLog.Debugf("activity received: \n%#v\n", activity)
			case <-tracker.quit:
				trackerLog.Infof("stopping activity tracker\n")
				//close all handlers for a clean exit
				for _, handler := range tracker.handlers {
					handler.Close()
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

//Start the tracker with all possible handlers
func (tracker *Instance) Start() (heartbeatCh chan *Heartbeat) {
	return tracker.StartWithHandlers(getAllHandlers()...)
}

//GetTrackerSystemState returns a copy of state of the system which is being used by the tracker
func (tracker *Instance) GetTrackerSystemState() system.State {
	return *tracker.state
}
