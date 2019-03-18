package service

import (
	"log"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

func MouseClickHandler(tickerCh chan struct{}, clickComm chan *activity.Type) (quit chan struct{}) {
	quit = make(chan struct{})
	registrationFree := make(chan struct{})

	go func(quit, registrationFree chan struct{}) {
		go addMouseClickRegistration(clickComm, registrationFree) //run once before first check
		for {
			select {
			case <-tickerCh:
				log.Printf("mouse clicker checked at : %v\n", time.Now())
				select {
				case _, ok := <-registrationFree:
					if ok {
						//log.Printf("registration free for mouse click \n")
						go addMouseClickRegistration(clickComm, registrationFree)
					} else {
						//log.Printf("error : channel closed \n")
						return
					}
				default:
					//nothing needs to be done
					//log.Printf("registration is busy for mouse click handler\n")
				}

			case <-quit:
				log.Printf("stopping click handler")
				return
			}
		}
	}(quit, registrationFree)
	return quit
}

func addMouseClickRegistration(clickComm chan *activity.Type, registrationFree chan struct{}) {
	log.Printf("adding reg \n")
	mleft := robotgo.AddEvent("mleft")
	if mleft {
		//log.Printf("mleft clicked \n")
		clickComm <- &activity.Type{
			ActivityType: activity.MOUSE_LEFT_CLICK,
		}
		registrationFree <- struct{}{}
		return
	}
}
