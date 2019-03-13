package activity

import (
	"log"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/prashantgupta24/activity-tracker/src/mouse"
)

func (tracker *ActivityTracker) Start() (heartbeatCh chan *Heartbeat, quit chan struct{}) {

	comm, quitMouseClickHandler := isMouseClicked(tracker)

	heartbeatCh = make(chan *Heartbeat, 1)
	quit = make(chan struct{})

	go func(tracker *ActivityTracker, heartbeatCh chan *Heartbeat, quit chan struct{}) {
		timeToCheck := tracker.TimeToCheck
		ticker := time.NewTicker(time.Second * timeToCheck)
		isIdle := true
		lastMousePos := mouse.GetPosition()
		for {
			select {
			case <-ticker.C:
				//log.Printf("tracker checking at %v\n", time.Now())
				currentMousePos := mouse.GetPosition()
				var heartbeat *Heartbeat
				if isIdle && isPointerIdle(currentMousePos, lastMousePos) {
					//log.Printf("no activity detected in the last %v seconds ...\n", int(timeToCheck))
					heartbeat = &Heartbeat{
						IsActivity: false,
						Time:       time.Now(),
					}
				} else {
					//log.Printf("activity detected in the last %v seconds ...\n", int(timeToCheck))
					heartbeat = &Heartbeat{
						IsActivity: true,
						Time:       time.Now(),
					}
					lastMousePos = currentMousePos
				}
				heartbeatCh <- heartbeat
				isIdle = true
			case <-comm:
				isIdle = false
				//log.Printf("value received: %v\n", isIdle)
			case <-quit:
				log.Printf("stopping activity tracker\n")
				quitMouseClickHandler <- struct{}{}
				//robotgo.StopEvent()
				return
			}
		}
	}(tracker, heartbeatCh, quit)

	return heartbeatCh, quit
}

func isPointerIdle(currentMousePos, lastMousePos *mouse.Position) bool {
	//log.Printf("current mouse position: %v\n", currentMousePos)
	//log.Printf("last mouse position: %v\n", lastMousePos)
	if currentMousePos.MouseX == lastMousePos.MouseX &&
		currentMousePos.MouseY == lastMousePos.MouseY {
		return true
	}
	return false
}

func isMouseClicked(tracker *ActivityTracker) (clickComm, quit chan struct{}) {
	ticker := time.NewTicker(time.Second * tracker.TimeToCheck)
	clickComm = make(chan struct{}, 1)
	quit = make(chan struct{})
	registrationFree := make(chan struct{})
	go func() {
		isRunning := false
		for {
			select {
			case <-ticker.C:
				//log.Printf("mouse clicker ticked at : %v\n", time.Now())
				if !isRunning {
					isRunning = true
					go func(registrationFree chan struct{}) {
						//log.Printf("adding reg \n")
						mleft := robotgo.AddEvent("mleft")
						if mleft {
							//log.Printf("mleft clicked \n")
							clickComm <- struct{}{}
							registrationFree <- struct{}{}
							return
						}
					}(registrationFree)
				}

				select {
				case _, ok := <-registrationFree:
					if ok {
						//log.Printf("registration free for mouse click \n")
						isRunning = false
					} else {
						//log.Printf("Channel closed \n")
					}
				default:
					//log.Printf("registration is busy for mouse click handler\n")
					isRunning = true
				}

			case <-quit:
				log.Printf("stopping click handler")
				close(clickComm)
				return
			}
		}
	}()
	return clickComm, quit
}
