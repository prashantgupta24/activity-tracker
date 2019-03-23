package service

import (
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
	log "github.com/sirupsen/logrus"
)

type testHandler struct {
	tickerCh chan struct{}
}

func (h *testHandler) Start(logger *log.Logger, activityCh chan *activity.Type) {
	h.tickerCh = make(chan struct{})
	go func(logger *log.Logger) {
		for range h.tickerCh {
			activityCh <- &activity.Type{
				ActivityType: activity.TEST_ACTIVITY,
			}
		}
		return
	}(logger)
}

func TestHandler() *testHandler {
	return &testHandler{}
}

func (h *testHandler) Type() activity.Type {
	return activity.Type{
		ActivityType: activity.TEST_ACTIVITY,
	}
}

func (h *testHandler) Trigger() {
	select {
	case h.tickerCh <- struct{}{}:
	default:
		//service is blocked, handle it somehow?
	}
}

func (h *testHandler) Close() {
	close(h.tickerCh)
}
