package tracker

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/handler"
	"github.com/prashantgupta24/activity-tracker/internal/pkg/logging"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

const (
	preHeartbeatTime = time.Millisecond * 100
	//heartbeat (seconds)
	minHFrequency     = 60
	maxHFrequency     = 300
	defaultHFrequency = 60

	//worker (seconds)
	minWFrequency     = 4
	maxWFrequency     = minHFrequency
	defaultWFrequency = minHFrequency

	numWorkerFrequencyDivisions = 5
)

//StartWithHandlers starts the tracker with a set of handlers
func (tracker *Instance) StartWithHandlers(handlers ...handler.Instance) (heartbeatCh chan *Heartbeat) {
	logger := logging.NewLoggerLevelFormat(tracker.LogLevel, tracker.LogFormat)

	//register handlers
	tracker.registerHandlers(logger, handlers...)

	//returned channel
	heartbeatCh = make(chan *Heartbeat, 1)

	tracker.quit = make(chan struct{})

	go func(logger *log.Logger, tracker *Instance) {
		trackerLog := logger.WithFields(log.Fields{
			"method": "activity-tracker",
		})

		//instantiating ticker frequencies
		heartbeatFrequency, workerFrequency := tracker.validateFrequencies()

		tickerHeartbeat := time.NewTicker(heartbeatFrequency * time.Second)
		tickerWorker := time.NewTicker(workerFrequency*time.Second - preHeartbeatTime)

		activityMap := makeActivityMap()

		for {
			select {
			case <-tickerWorker.C:
				trackerLog.Debugln("tracker worker working")
				//time to trigger all registered handlers
				for _, handler := range tracker.handlers {
					handler.Trigger()
				}
			case <-tickerHeartbeat.C:
				trackerLog.Debugln("tracker heartbeat checking")
				var heartbeat *Heartbeat
				if len(activityMap) == 0 {
					logger.Debugf("no activity detected in the last %v seconds ...\n", int(heartbeatFrequency))
					heartbeat = &Heartbeat{
						WasAnyActivity: false,
						ActivityMap:    nil,
						Time:           time.Now(),
					}
				} else {
					trackerLog.Debugf("activity detected in the last %v seconds ...\n", int(heartbeatFrequency))
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

func getAllHandlers() []handler.Instance {
	return []handler.Instance{
		handler.MouseClickHandler(), handler.MouseCursorHandler(),
		handler.ScreenChangeHandler(),
	}
}

func makeActivityMap() (activityMap map[activity.Type][]time.Time) {
	activityMap = make(map[activity.Type][]time.Time)
	return activityMap
}

func (tracker *Instance) validateFrequencies() (heartbeatFreqReturn, workerFreqReturn time.Duration) {
	heartbeatFreq := tracker.HeartbeatFrequency
	workerFreq := tracker.WorkerFrequency

	if tracker.isTest {
		heartbeatFreqReturn = time.Duration(heartbeatFreq)
		workerFreqReturn = time.Duration(heartbeatFreq)
		return
	}

	//heartbeat check
	if heartbeatFreq >= minHFrequency && heartbeatFreq <= maxHFrequency {
		heartbeatFreqReturn = time.Duration(heartbeatFreq) //within range
	} else {
		heartbeatFreqReturn = time.Duration(defaultHFrequency)
	}

	//worker check
	if workerFreq >= minWFrequency && workerFreq <= maxWFrequency {
		workerFreqReturn = time.Duration(workerFreq) //within range
	} else {
		workerFreqReturn = time.Duration(defaultWFrequency)
	}
	return
}

func (tracker *Instance) registerHandlers(logger *log.Logger, handlers ...handler.Instance) {

	tracker.handlers = make(map[activity.Type]handler.Instance)
	tracker.activityCh = make(chan *activity.Instance, len(handlers)) // number based on types of activities being tracked

	for _, handler := range handlers {
		if _, ok := tracker.handlers[handler.Type()]; !ok { //duplicate registration prevention
			tracker.handlers[handler.Type()] = handler
			handler.Start(logger, tracker.activityCh)
		}
	}
}
