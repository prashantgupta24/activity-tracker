package handler

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-vgo/robotgo"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
	"github.com/prashantgupta24/activity-tracker/pkg/system"
)

const (
	screenChangeActivity = activity.ScreenChange
)

//ScreenChangeHandlerStruct is the handler for screen changes
type ScreenChangeHandlerStruct struct {
	screenHandlerLogger *log.Entry
	tickerCh            chan struct{}
}

type screenInfo struct {
	didScreenChange    bool
	currentScreenTitle string
}

//Start the handler
func (s *ScreenChangeHandlerStruct) Start(logger *log.Logger, activityCh chan *activity.Instance) {

	s.tickerCh = make(chan struct{})
	s.screenHandlerLogger = logger.WithFields(log.Fields{
		"method": "screen-change-handler",
	})

	go func() {
		lastScreenTitle := robotgo.GetTitle()

		for range s.tickerCh {
			s.screenHandlerLogger.Debugf("screen change checked")
			commCh := make(chan *screenInfo)
			go checkScreenChange(s.screenHandlerLogger, commCh, lastScreenTitle)
			select {
			case screenInfo := <-commCh:
				if screenInfo.didScreenChange {
					activityCh <- &activity.Instance{
						Type: screenChangeActivity,
					}
					lastScreenTitle = screenInfo.currentScreenTitle
				}
			case <-time.After(timeout * time.Millisecond):
				//timeout, do nothing
				s.screenHandlerLogger.Debugf("timeout happened after %vms while checking screen change handler", timeout)
			}
		}
		s.screenHandlerLogger.Infof("stopping screen change handler")
		return
	}()
}

//ScreenChangeHandler returns an instance of the struct
func ScreenChangeHandler() *ScreenChangeHandlerStruct {
	return &ScreenChangeHandlerStruct{}
}

//Trigger the handler
func (s *ScreenChangeHandlerStruct) Trigger(state system.State) {
	//no point triggering the handler since the system is asleep
	if state.IsSystemSleep {
		s.screenHandlerLogger.Debugf("system sleeping so not working")
		return
	}
	//doing it the non-blocking sender way
	select {
	case s.tickerCh <- struct{}{}:
	default:
		//handler is blocked, handle it somehow?
	}
}

//Type returns the type of handler
func (s *ScreenChangeHandlerStruct) Type() activity.Type {
	return screenChangeActivity
}

//Close closes the handler
func (s *ScreenChangeHandlerStruct) Close() {
	close(s.tickerCh)
}

func checkScreenChange(logger *log.Entry, commCh chan *screenInfo, lastScreenTitle string) {
	currentScreenTitle := robotgo.GetTitle()
	screenLogger := logger.WithFields(log.Fields{
		"current-title": currentScreenTitle,
		"last-title":    lastScreenTitle,
	})
	if lastScreenTitle != currentScreenTitle {
		commCh <- &screenInfo{
			didScreenChange:    true,
			currentScreenTitle: currentScreenTitle,
		}
		screenLogger.Debugf("screen changed")
	} else { //comment this section out to test timeout logic
		commCh <- &screenInfo{
			didScreenChange:    false,
			currentScreenTitle: "",
		}
		screenLogger.Debugf("screen did not change")
	}
}
