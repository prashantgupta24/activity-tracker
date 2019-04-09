package handler

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/mouse"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
	"github.com/prashantgupta24/activity-tracker/pkg/system"
)

const (
	mouseMoveActivity = activity.MouseCursorMovement
)

//MouseCursorHandlerStruct is the handler for mouse cursor movements
type MouseCursorHandlerStruct struct {
	cursurHandlerLogger *log.Entry
	tickerCh            chan struct{}
}

type cursorInfo struct {
	didCursorMove   bool
	currentMousePos *mouse.Position
}

//Start the handler
func (m *MouseCursorHandlerStruct) Start(logger *log.Logger, activityCh chan *activity.Instance) {

	m.tickerCh = make(chan struct{})
	m.cursurHandlerLogger = logger.WithFields(log.Fields{
		"method": "mouse-cursor-handler",
	})

	go func() {
		lastMousePos := mouse.GetPosition()
		for range m.tickerCh {
			m.cursurHandlerLogger.Debugf("mouse cursor checked")
			commCh := make(chan *cursorInfo)
			go checkCursorChange(m.cursurHandlerLogger, commCh, lastMousePos)
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
				m.cursurHandlerLogger.Debugf("timeout happened after %vms while checking mouse cursor handler", timeout)
			}
		}
		m.cursurHandlerLogger.Infof("stopping cursor handler")
		return
	}()
}

//MouseCursorHandler returns an instance of the struct
func MouseCursorHandler() *MouseCursorHandlerStruct {
	return &MouseCursorHandlerStruct{}
}

//Trigger the handler
func (m *MouseCursorHandlerStruct) Trigger(state system.State) {
	//no point triggering the handler since the system is asleep
	if state.IsSystemSleep {
		m.cursurHandlerLogger.Debugf("system sleeping so not working")
		return
	}
	//doing it the non-blocking sender way
	select {
	case m.tickerCh <- struct{}{}:
	default:
		//handler is blocked, handle it somehow?
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
