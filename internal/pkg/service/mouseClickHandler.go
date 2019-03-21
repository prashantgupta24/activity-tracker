package service

import (
	"time"

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

	go func() {
		go addMouseClickRegistration(logger, activityCh, registrationFree) //run once before first check
		for range m.tickerCh {
			logger.Debugf("mouse clicker checked at : %v\n", time.Now())
			select {
			case _, ok := <-registrationFree:
				if ok {
					logger.Debugf("registration free for mouse click \n")
					go addMouseClickRegistration(logger, activityCh, registrationFree)
				} else {
					logger.Errorf("error : channel closed \n")
					return
				}
			default:
				logger.Debugf("registration is busy for mouse click handler, do nothing\n")
			}
		}
		logger.Infof("stopping click handler")
		return
	}()
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
func (m *mouseClickHandler) Close() {
	close(m.tickerCh)
}

func addMouseClickRegistration(logger *log.Logger, activityCh chan *activity.Type,
	registrationFree chan struct{}) {

	logger.Debugf("adding reg \n")
	mleft := robotgo.AddEvent("mleft")
	if mleft {
		logger.Debugf("mleft clicked \n")
		activityCh <- &activity.Type{
			ActivityType: activity.MOUSE_LEFT_CLICK,
		}
		registrationFree <- struct{}{}
		return
	}
}
