package service

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/mouse"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

type mouseCursorHandler struct {
	tickerCh chan struct{}
}

type cursorInfo struct {
	didCursorMove   bool
	currentMousePos *mouse.Position
}

func (m *mouseCursorHandler) Start(logger *log.Logger, activityCh chan *activity.Type) {

	m.tickerCh = make(chan struct{})

	go func() {
		lastMousePos := mouse.GetPosition()
		for range m.tickerCh {
			logger.Debugf("mouse cursor checked at : %v\n", time.Now())
			commCh := make(chan *cursorInfo)
			go checkCursorChange(logger, commCh, lastMousePos)
			select {
			case cursorInfo := <-commCh:
				if cursorInfo.didCursorMove {
					activityCh <- &activity.Type{
						ActivityType: activity.MOUSE_CURSOR_MOVEMENT,
					}
					lastMousePos = cursorInfo.currentMousePos
				}
			case <-time.After(timeout * time.Millisecond):
				//timeout, do nothing
				logger.Debugf("timeout happened after %vms while checking mouse cursor handler", timeout)
			}
		}
		logger.Infof("stopping cursor handler")
		return
	}()
}

func MouseCursorHandler() *mouseCursorHandler {
	return &mouseCursorHandler{}
}

func (m *mouseCursorHandler) Trigger() {
	//doing it the non-blocking sender way
	select {
	case m.tickerCh <- struct{}{}:
	default:
		//service is blocked, handle it somehow?
	}
}
func (m *mouseCursorHandler) Close() {
	close(m.tickerCh)
}

func checkCursorChange(logger *log.Logger, commCh chan *cursorInfo, lastMousePos *mouse.Position) {
	currentMousePos := mouse.GetPosition()
	logger.Debugf("current mouse position: %v\n", currentMousePos)
	logger.Debugf("last mouse position: %v\n", lastMousePos)
	if currentMousePos.MouseX == lastMousePos.MouseX &&
		currentMousePos.MouseY == lastMousePos.MouseY {
		commCh <- &cursorInfo{
			didCursorMove:   false,
			currentMousePos: nil,
		}
	} else {
		commCh <- &cursorInfo{
			didCursorMove:   true,
			currentMousePos: currentMousePos,
		}
	}
}
