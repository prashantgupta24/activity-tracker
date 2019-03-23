package service

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/mouse"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

const (
	mouseMoveActivity = activity.MOUSE_CURSOR_MOVEMENT
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

	go func(logger *log.Logger) {
		handlerLogger := logger.WithFields(log.Fields{
			"method": "mouse-cursor-handler",
		})
		lastMousePos := mouse.GetPosition()
		for range m.tickerCh {
			handlerLogger.Debugf("mouse cursor checked")
			commCh := make(chan *cursorInfo)
			go checkCursorChange(handlerLogger, commCh, lastMousePos)
			select {
			case cursorInfo := <-commCh:
				if cursorInfo.didCursorMove {
					activityCh <- &activity.Type{
						ActivityType: mouseMoveActivity,
					}
					lastMousePos = cursorInfo.currentMousePos
				}
			case <-time.After(timeout * time.Millisecond):
				//timeout, do nothing
				handlerLogger.Debugf("timeout happened after %vms while checking mouse cursor handler", timeout)
			}
		}
		handlerLogger.Infof("stopping cursor handler")
		return
	}(logger)
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

func (m *mouseCursorHandler) Type() activity.Type {
	return activity.Type{
		ActivityType: mouseMoveActivity,
	}
}

func (m *mouseCursorHandler) Close() {
	close(m.tickerCh)
}

func checkCursorChange(logger *log.Entry, commCh chan *cursorInfo, lastMousePos *mouse.Position) {
	currentMousePos := mouse.GetPosition()
	cursorLogger := logger.WithFields(log.Fields{
		"current-mouse-position": currentMousePos,
		"last-mouse-position":    lastMousePos,
	})
	if currentMousePos.MouseX == lastMousePos.MouseX &&
		currentMousePos.MouseY == lastMousePos.MouseY {
		commCh <- &cursorInfo{
			didCursorMove:   false,
			currentMousePos: nil,
		}
		cursorLogger.Debugf("cursor not moved")
	} else {
		commCh <- &cursorInfo{
			didCursorMove:   true,
			currentMousePos: currentMousePos,
		}
		cursorLogger.Debugf("cursor moved")
	}
}
