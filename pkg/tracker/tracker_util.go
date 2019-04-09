package tracker

import (
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/handler"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
	log "github.com/sirupsen/logrus"
)

//add handlers to this function to start with them by default
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

func (tracker *Instance) validateIntervals() (heartbeatIntervalReturn, workerIntervalReturn time.Duration) {
	heartbeatInterval := tracker.HeartbeatInterval
	workerInterval := tracker.WorkerInterval

	if tracker.isTest {
		heartbeatIntervalReturn = time.Duration(heartbeatInterval)
		workerIntervalReturn = time.Duration(heartbeatInterval)
		return
	}

	//heartbeat check
	if heartbeatInterval >= minHInterval && heartbeatInterval <= maxHInterval {
		heartbeatIntervalReturn = time.Duration(heartbeatInterval) //within range
	} else {
		heartbeatIntervalReturn = time.Duration(defaultHInterval)
	}

	//worker check
	if workerInterval >= minWInterval && workerInterval <= maxWInterval {
		workerIntervalReturn = time.Duration(workerInterval) //within range
	} else {
		workerIntervalReturn = time.Duration(defaultWInterval)
	}
	return
}

//validate all handlers with certain rules
func validateHandlers(handlers ...handler.Instance) []handler.Instance {

	isMachineSleepHandlerPresent := false

	for _, handler := range handlers {
		switch handler.Type() {
		case activity.MachineSleep, activity.MachineWake:
			isMachineSleepHandlerPresent = true
		}
	}

	//condition 1, adding machine sleep handler as a fail-safe in all scenarios
	if !isMachineSleepHandlerPresent {
		handlers = append(handlers, handler.MachineSleepHandler())
	}
	return handlers
}

func (tracker *Instance) registerHandlers(logger *log.Logger, handlers ...handler.Instance) {
	//validate handlers first
	if !tracker.isTest {
		handlers = validateHandlers(handlers...)
	}

	tracker.handlers = make(map[activity.Type]handler.Instance)
	tracker.activityCh = make(chan *activity.Instance, len(handlers)) // number based on types of activities being tracked

	for _, handler := range handlers {
		if _, ok := tracker.handlers[handler.Type()]; !ok { //duplicate registration prevention
			tracker.handlers[handler.Type()] = handler
			handler.Start(logger, tracker.activityCh)
		}
	}
}
