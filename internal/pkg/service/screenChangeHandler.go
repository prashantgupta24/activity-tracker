package service

import (
	"log"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

func ScreenChangeHandler(comm chan *activity.Type) (tickerCh chan struct{}) {

	tickerCh = make(chan struct{})

	go func() {
		screenSizeX, screenSizeY := robotgo.GetScreenSize()
		pixelPointX := int(screenSizeX / 2)
		pixelPointY := int(screenSizeY / 2)
		lastPixelColor := robotgo.GetPixelColor(pixelPointX, pixelPointY)
		for range tickerCh {
			log.Printf("screen change checked at : %v\n", time.Now())
			currentPixelColor := robotgo.GetPixelColor(pixelPointX, pixelPointY)
			// log.Printf("current pixel color: %v\n", currentPixelColor)
			// log.Printf("last pixel color: %v\n", lastPixelColor)
			if lastPixelColor != currentPixelColor {
				comm <- &activity.Type{
					ActivityType: activity.SCREEN_CHANGE,
				}
				lastPixelColor = currentPixelColor
			}
		}
		log.Printf("stopping screen change handler")
		return
	}()

	return tickerCh
}
