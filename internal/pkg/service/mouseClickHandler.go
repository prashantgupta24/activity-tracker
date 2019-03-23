package service

import (
	log "github.com/sirupsen/logrus"

	"github.com/go-vgo/robotgo"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

type mouseClickHandler struct {
	tickerCh chan struct{}
}

func (m *mouseClickHandler) Start(logger *log.Logger, activityCh chan *activity.Type) {
	m.tickerCh = make(chan struct{})
	registrationFree := make(chan struct{})

	go func(logger *log.Logger) {
		handlerLogger := logger.WithFields(log.Fields{
			"method": "mouse-click-handler",
		})
		go addMouseClickRegistration(handlerLogger, activityCh, registrationFree) //run once before first check
		for range m.tickerCh {
			handlerLogger.Debugln("mouse clicker checked")
			select {
			case _, ok := <-registrationFree:
				if ok {
					handlerLogger.Debugf("registration free \n")
					go addMouseClickRegistration(handlerLogger, activityCh, registrationFree)
				} else {
					handlerLogger.Errorf("error : channel closed \n")
					return
				}
			default:
				handlerLogger.Debugf("registration is busy, do nothing\n")
			}
		}
		handlerLogger.Infof("stopping click handler")
		return
	}(logger)
}

func MouseClickHandler() *mouseClickHandler {
	return &mouseClickHandler{}
}

func (m *mouseClickHandler) Trigger() {
	//doing it the non-blocking sender way
	select {
	case m.tickerCh <- struct{}{}:
	default:
		//service is blocked, handle it somehow?
	}
}

func (m *mouseClickHandler) Type() activity.Type {
	return activity.Type{
		ActivityType: activity.MOUSE_CLICK,
	}
}

func (m *mouseClickHandler) Close() {
	close(m.tickerCh)
}

func addMouseClickRegistration(logger *log.Entry, activityCh chan *activity.Type,
	registrationFree chan struct{}) {

	logger.Debugf("adding mouse left click registration \n")
	mleft := robotgo.AddEvent("mleft")
	if mleft {
		logger.Debugf("mleft clicked \n")
		activityCh <- &activity.Type{
			ActivityType: activity.MOUSE_CLICK,
		}
		registrationFree <- struct{}{}
		return
	}
}
