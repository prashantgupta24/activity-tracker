package service

import (
	"log"
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/mouse"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

func MouseCursorHandler(tickerCh chan struct{}, comm chan *activity.Type) (quit chan struct{}) {

	quit = make(chan struct{})

	go func(quit chan struct{}) {
		lastMousePos := mouse.GetPosition()
		for {
			select {
			case <-tickerCh:
				log.Printf("mouse cursor checked at : %v\n", time.Now())
				//log.Printf("current mouse position: %v\n", currentMousePos)
				//log.Printf("last mouse position: %v\n", lastMousePos)
				currentMousePos := mouse.GetPosition()
				if currentMousePos.MouseX == lastMousePos.MouseX &&
					currentMousePos.MouseY == lastMousePos.MouseY {
					continue
				}
				comm <- &activity.Type{
					ActivityType: activity.MOUSE_CURSOR_MOVEMENT,
				}
				lastMousePos = currentMousePos
			case <-quit:
				log.Printf("stopping cursor handler")
				return
			}
		}
	}(quit)
	return quit
}
