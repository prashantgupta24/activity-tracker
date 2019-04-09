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
	tickerCh chan struct{}
}

type screenInfo struct {
	didScreenChange   bool
	currentPixelColor string
}

//Start the handler
func (s *ScreenChangeHandlerStruct) Start(logger *log.Logger, activityCh chan *activity.Instance) {

	s.tickerCh = make(chan struct{})

	go func(logger *log.Logger) {
		handlerLogger := logger.WithFields(log.Fields{
			"method": "screen-change-handler",
		})
		screenSizeX, screenSizeY := robotgo.GetScreenSize()
		pixelPointX := int(screenSizeX / 2)
		pixelPointY := int(screenSizeY / 2)
		lastPixelColor := robotgo.GetPixelColor(pixelPointX, pixelPointY)
		for range s.tickerCh {
			handlerLogger.Debugf("screen change checked")
			commCh := make(chan *screenInfo)
			go checkScreenChange(handlerLogger, commCh, lastPixelColor, pixelPointX, pixelPointY)
			select {
			case screenInfo := <-commCh:
				if screenInfo.didScreenChange {
					activityCh <- &activity.Instance{
						Type: screenChangeActivity,
					}
					lastPixelColor = screenInfo.currentPixelColor
				}
			case <-time.After(timeout * time.Millisecond):
				//timeout, do nothing
				handlerLogger.Debugf("timeout happened after %vms while checking screen change handler", timeout)
			}
		}
		handlerLogger.Infof("stopping screen change handler")
		return
	}(logger)
}

//ScreenChangeHandler returns an instance of the struct
func ScreenChangeHandler() *ScreenChangeHandlerStruct {
	return &ScreenChangeHandlerStruct{}
}

//Trigger the handler
func (s *ScreenChangeHandlerStruct) Trigger(state system.State) {
	//no point triggering the handler since the system is asleep
	if state.IsSystemSleep {
		log.Infof("%v system sleeping so screen checker not doing any checking", time.Now())
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

func checkScreenChange(logger *log.Entry, commCh chan *screenInfo, lastPixelColor string, pixelPointX, pixelPointY int) {
	currentPixelColor := robotgo.GetPixelColor(pixelPointX, pixelPointY)
	screenLogger := logger.WithFields(log.Fields{
		"current-pixel": currentPixelColor,
		"last-pixel":    lastPixelColor,
	})
	//robotgo.MoveMouse(pixelPointX, pixelPointY)
	if lastPixelColor != currentPixelColor {
		commCh <- &screenInfo{
			didScreenChange:   true,
			currentPixelColor: currentPixelColor,
		}
		screenLogger.Debugf("screen changed")
	} else { //comment this section out to test timeout logic
		commCh <- &screenInfo{
			didScreenChange:   false,
			currentPixelColor: "",
		}
		screenLogger.Debugf("screen did not change")
	}
}
