package service

import (
	"log"
	"time"

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

func (s *screenChangeHandler) Start(activityCh chan *activity.Type) {

	s.tickerCh = make(chan struct{})

	go func() {
		screenSizeX, screenSizeY := robotgo.GetScreenSize()
		pixelPointX := int(screenSizeX / 2)
		pixelPointY := int(screenSizeY / 2)
		lastPixelColor := robotgo.GetPixelColor(pixelPointX, pixelPointY)
		for range s.tickerCh {
			log.Printf("screen change checked at : %v\n", time.Now())
			commCh := make(chan *screenInfo)
			go checkScreenChange(commCh, lastPixelColor, pixelPointX, pixelPointY)
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
				log.Printf("timeout happened after %vms while checking screen change handler", timeout)
			}
		}
		log.Printf("stopping screen change handler")
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

func checkScreenChange(commCh chan *screenInfo, lastPixelColor string, pixelPointX, pixelPointY int) {
	currentPixelColor := robotgo.GetPixelColor(pixelPointX, pixelPointY)
	// log.Printf("current pixel color: %v\n", currentPixelColor)
	// log.Printf("last pixel color: %v\n", lastPixelColor)
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
