package service

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-vgo/robotgo"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

type screenChangeHandler struct {
	tickerCh chan struct{}
}

type screenInfo struct {
	didScreenChange   bool
	currentPixelColor string
}

func (s *screenChangeHandler) Start(logger *log.Logger, activityCh chan *activity.Type) {

	s.tickerCh = make(chan struct{})

	go func() {
		screenSizeX, screenSizeY := robotgo.GetScreenSize()
		pixelPointX := int(screenSizeX / 2)
		pixelPointY := int(screenSizeY / 2)
		lastPixelColor := robotgo.GetPixelColor(pixelPointX, pixelPointY)
		for range s.tickerCh {
			logger.Debugf("screen change checked at : %v\n", time.Now())
			commCh := make(chan *screenInfo)
			go checkScreenChange(logger, commCh, lastPixelColor, pixelPointX, pixelPointY)
			select {
			case screenInfo := <-commCh:
				if screenInfo.didScreenChange {
					activityCh <- &activity.Type{
						ActivityType: activity.SCREEN_CHANGE,
					}
					lastPixelColor = screenInfo.currentPixelColor
				}
			case <-time.After(timeout * time.Millisecond):
				//timeout, do nothing
				logger.Debugf("timeout happened after %vms while checking screen change handler", timeout)
			}
		}
		logger.Infof("stopping screen change handler")
		return
	}()
}

func ScreenChangeHandler() *screenChangeHandler {
	return &screenChangeHandler{}
}

func (s *screenChangeHandler) Trigger() {
	//doing it the non-blocking sender way
	select {
	case s.tickerCh <- struct{}{}:
	default:
		//service is blocked, handle it somehow?
	}
}
func (s *screenChangeHandler) Close() {
	close(s.tickerCh)
}

func checkScreenChange(logger *log.Logger, commCh chan *screenInfo, lastPixelColor string, pixelPointX, pixelPointY int) {
	currentPixelColor := robotgo.GetPixelColor(pixelPointX, pixelPointY)
	logger.Debugf("current pixel color: %v\n", currentPixelColor)
	logger.Debugf("last pixel color: %v\n", lastPixelColor)
	//robotgo.MoveMouse(pixelPointX, pixelPointY)
	if lastPixelColor != currentPixelColor {
		commCh <- &screenInfo{
			didScreenChange:   true,
			currentPixelColor: currentPixelColor,
		}
	} else { //comment this section out to test timeout logic
		commCh <- &screenInfo{
			didScreenChange:   false,
			currentPixelColor: "",
		}
	}
}
