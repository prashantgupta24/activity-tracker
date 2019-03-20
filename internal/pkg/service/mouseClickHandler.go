package service

import (
	"log"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

type mouseClickHandler struct {
	tickerCh chan struct{}
}

func (m *mouseClickHandler) Start(activityCh chan *activity.Type) {
	m.tickerCh = make(chan struct{})
	registrationFree := make(chan struct{})

	go func() {
		go addMouseClickRegistration(activityCh, registrationFree) //run once before first check
		for range m.tickerCh {
			log.Printf("mouse clicker checked at : %v\n", time.Now())
			select {
			case _, ok := <-registrationFree:
				if ok {
					//log.Printf("registration free for mouse click \n")
					go addMouseClickRegistration(activityCh, registrationFree)
				} else {
					//log.Printf("error : channel closed \n")
					return
				}
			default:
				//log.Printf("registration is busy for mouse click handler, do nothing\n")
			}
		}
		log.Printf("stopping click handler")
		return
	}()
}

func MouseClickHandler() *mouseClickHandler {
	return &mouseClickHandler{}
}

func (m *mouseClickHandler) Trigger() {
	//doing it the non-blocking sender way
	select {
	case m.tickerCh <- struct{}{}:
	default:
		//service is blocked, handle it somehow?
	}
}
func (m *mouseClickHandler) Close() {
	close(m.tickerCh)
}

func addMouseClickRegistration(activityCh chan *activity.Type, registrationFree chan struct{}) {
	log.Printf("adding reg \n")
	mleft := robotgo.AddEvent("mleft")
	if mleft {
		//log.Printf("mleft clicked \n")
		activityCh <- &activity.Type{
			ActivityType: activity.MOUSE_LEFT_CLICK,
		}
		registrationFree <- struct{}{}
		return
	}
}
