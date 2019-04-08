package handler

import (
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
	"github.com/prashantgupta24/activity-tracker/pkg/system"
	log "github.com/sirupsen/logrus"
)

const (
	testActivity = activity.TestActivity
)

//TestHandlerStruct is a test handler
type TestHandlerStruct struct {
	testHandlerLogger *log.Entry
	tickerCh          chan struct{}
}

//Start the handler
func (h *TestHandlerStruct) Start(logger *log.Logger, activityCh chan *activity.Instance) {
	h.tickerCh = make(chan struct{})
	h.testHandlerLogger = logger.WithFields(log.Fields{
		"method": "test-handler",
	})

	go func() {
		for range h.tickerCh {
			activityCh <- &activity.Instance{
				Type: testActivity,
			}
		}
		h.testHandlerLogger.Infof("stopping test handler")
		return
	}()
}

//TestHandler returns an instance of the struct
func TestHandler() *TestHandlerStruct {
	return &TestHandlerStruct{}
}

//Type returns the type of handler
func (h *TestHandlerStruct) Type() activity.Type {
	return testActivity
}

//Trigger the handler
func (h *TestHandlerStruct) Trigger(state system.State) {
	//no point triggering the handler since the system is asleep
	if state.IsSystemSleep {
		h.testHandlerLogger.Debugf("system sleeping so not working")
		return
	}
	select {
	case h.tickerCh <- struct{}{}:
	default:
		//handler is blocked, handle it somehow?
	}
}

//Close closes the handler
func (h *TestHandlerStruct) Close() {
	close(h.tickerCh)
}
