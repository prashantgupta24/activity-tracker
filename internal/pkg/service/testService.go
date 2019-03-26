package service

import (
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
	log "github.com/sirupsen/logrus"
)

const (
	testActivity = activity.TestActivity
)

//TestHandlerStruct is a test handler
type TestHandlerStruct struct {
	tickerCh chan struct{}
}

//Start the service
func (h *TestHandlerStruct) Start(logger *log.Logger, activityCh chan *activity.Type) {
	h.tickerCh = make(chan struct{})
	go func() {
		for range h.tickerCh {
			activityCh <- &activity.Type{
				ActivityType: testActivity,
			}
		}
		logger.Infof("stopping test handler")
		return
	}()
}

//TestHandler returns an instance of the struct
func TestHandler() *TestHandlerStruct {
	return &TestHandlerStruct{}
}

//Type return the type of handler
func (h *TestHandlerStruct) Type() activity.Type {
	return activity.Type{
		ActivityType: testActivity,
	}
}

//Trigger the service
func (h *TestHandlerStruct) Trigger() {
	select {
	case h.tickerCh <- struct{}{}:
	default:
		//service is blocked, handle it somehow?
	}
}

//Close closes the handler
func (h *TestHandlerStruct) Close() {
	close(h.tickerCh)
}
