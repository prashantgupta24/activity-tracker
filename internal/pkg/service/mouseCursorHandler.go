package service

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/mouse"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

const (
	mouseMoveActivity = activity.MouseCursorMovement
)

//MouseCursorHandlerStruct is the handler for mouse cursor movements
type MouseCursorHandlerStruct struct {
	tickerCh chan struct{}
}

type cursorInfo struct {
	didCursorMove   bool
	currentMousePos *mouse.Position
}

//Start the service
func (m *MouseCursorHandlerStruct) Start(logger *log.Logger, activityCh chan *activity.Instance) {

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
					activityCh <- &activity.Instance{
						Type: mouseMoveActivity,
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

//MouseCursorHandler returns an instance of the struct
func MouseCursorHandler() *MouseCursorHandlerStruct {
	return &MouseCursorHandlerStruct{}
}

//Trigger the service
func (m *MouseCursorHandlerStruct) Trigger() {
	//doing it the non-blocking sender way
	select {
	case m.tickerCh <- struct{}{}:
	default:
		//service is blocked, handle it somehow?
	}
}

//Type returns the type of handler
func (m *MouseCursorHandlerStruct) Type() activity.Type {
	return mouseMoveActivity
}

//Close closes the handler
func (m *MouseCursorHandlerStruct) Close() {
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
