package service

import (
	"log"
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/mouse"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

func MouseCursorHandler(comm chan *activity.Type) (tickerCh chan struct{}) {

	tickerCh = make(chan struct{})

	go func() {
		lastMousePos := mouse.GetPosition()
		for range tickerCh {
			log.Printf("mouse cursor checked at : %v\n", time.Now())
			currentMousePos := mouse.GetPosition()
			//log.Printf("current mouse position: %v\n", currentMousePos)
			//log.Printf("last mouse position: %v\n", lastMousePos)
			if currentMousePos.MouseX == lastMousePos.MouseX &&
				currentMousePos.MouseY == lastMousePos.MouseY {
				continue
			}
			comm <- &activity.Type{
				ActivityType: activity.MOUSE_CURSOR_MOVEMENT,
			}
			lastMousePos = currentMousePos
		}
		log.Printf("stopping cursor handler")
		return
	}()

	return tickerCh
}
