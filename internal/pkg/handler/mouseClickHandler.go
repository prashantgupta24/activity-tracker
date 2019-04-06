package handler

import (
	log "github.com/sirupsen/logrus"

	"github.com/go-vgo/robotgo"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
	"github.com/prashantgupta24/activity-tracker/pkg/system"
)

const (
	mouseClickActivity = activity.MouseClick
)

//MouseClickHandlerStruct is the handler for mouse clicks
type MouseClickHandlerStruct struct {
	tickerCh chan struct{}
}

//Start the handler
func (m *MouseClickHandlerStruct) Start(logger *log.Logger, activityCh chan *activity.Instance) {
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

//MouseClickHandler returns an instance of the struct
func MouseClickHandler() *MouseClickHandlerStruct {
	return &MouseClickHandlerStruct{}
}

//Trigger the handler
func (m *MouseClickHandlerStruct) Trigger(state system.State) {
	//no point triggering the handler since the system is asleep
	if state.IsSystemSleep {
		return
	}
	//doing it the non-blocking sender way
	select {
	case m.tickerCh <- struct{}{}:
	default:
		//handler is blocked, handle it somehow?
	}
}

//Type returns the type of handler
func (m *MouseClickHandlerStruct) Type() activity.Type {
	return mouseClickActivity
}

//Close closes the handler
func (m *MouseClickHandlerStruct) Close() {
	close(m.tickerCh)
}

func addMouseClickRegistration(logger *log.Entry, activityCh chan *activity.Instance,
	registrationFree chan struct{}) {

	logger.Debugf("adding mouse left click registration \n")
	mleft := robotgo.AddEvent("mleft")
	if mleft {
		logger.Debugf("mleft clicked \n")
		activityCh <- &activity.Instance{
			Type: mouseClickActivity,
		}
		registrationFree <- struct{}{}
		return
	}
}
