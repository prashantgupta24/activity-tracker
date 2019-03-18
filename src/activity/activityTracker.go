package activity

import (
	"log"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/prashantgupta24/activity-tracker/src/mouse"
)

const (
	preHeartbeatTime = time.Millisecond * 10
)

func (tracker *ActivityTracker) Start() (heartbeatCh chan *Heartbeat, quit chan struct{}) {

	tracker.registerServices(handleMouseClicks, checkMouseCursorMovement)
	//returned channels
	heartbeatCh = make(chan *Heartbeat, 1)
	quit = make(chan struct{})

	go func(tracker *ActivityTracker, heartbeatCh chan *Heartbeat, quit chan struct{}) {

		timeToCheck := tracker.TimeToCheck
		//tickers
		tickerHeartbeat := time.NewTicker(time.Second * timeToCheck)
		tickerWorker := time.NewTicker(time.Second*timeToCheck - preHeartbeatTime)

		activities := makeActivityMap()

		//pixel check
		//screenSizeX, screenSizeY := robotgo.GetScreenSize()
		//lastPixelColor := robotgo.GetPixelColor(int(screenSizeX/2), int(screenSizeY/2))

		for {
			select {
			case <-tickerWorker.C:
				log.Printf("tracker worker working at %v\n", time.Now())
				for i := 0; i < len(tracker.services); i++ {
					tracker.workerTickerCh <- struct{}{}
				}
			case <-tickerHeartbeat.C:
				log.Printf("tracker heartbeat checking at %v\n", time.Now())
				var heartbeat *Heartbeat
				if len(activities) == 0 {
					//log.Printf("no activity detected in the last %v seconds ...\n", int(timeToCheck))
					heartbeat = &Heartbeat{
						IsActivity: false,
						Activity:   nil,
						Time:       time.Now(),
					}
				} else {
					//log.Printf("activity detected in the last %v seconds ...\n", int(timeToCheck))
					heartbeat = &Heartbeat{
						IsActivity: true,
						Activity:   activities,
						Time:       time.Now(),
					}

				}
				heartbeatCh <- heartbeat
				activities = makeActivityMap() //reset the activities map
			case activity := <-tracker.activityCh:
				activities[activity] = time.Now()
				//log.Printf("activity received: %#v\n", activity)
			case <-quit:
				log.Printf("stopping activity tracker\n")
				for _, quitHandler := range tracker.services {
					quitHandler <- struct{}{}
				}
				//robotgo.StopEvent()
				return
			}
		}
	}(tracker, heartbeatCh, quit)

	return heartbeatCh, quit
}

func checkMouseCursorMovement(tickerCh chan struct{}, comm chan *Activity) (quit chan struct{}) {

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
				comm <- &Activity{
					ActivityType: MOUSE_CURSOR_MOVEMENT,
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

func addMouseClickRegistration(clickComm chan *Activity, registrationFree chan struct{}) {
	log.Printf("adding reg \n")
	mleft := robotgo.AddEvent("mleft")
	if mleft {
		//log.Printf("mleft clicked \n")
		clickComm <- &Activity{
			ActivityType: MOUSE_LEFT_CLICK,
		}
		registrationFree <- struct{}{}
		return
	}
}

func handleMouseClicks(tickerCh chan struct{}, clickComm chan *Activity) (quit chan struct{}) {
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
						//isRunning = false
						go addMouseClickRegistration(clickComm, registrationFree)
					} else {
						//log.Printf("error : channel closed \n")
						return
					}
				default:
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

func makeActivityMap() map[*Activity]time.Time {
	activityMap := make(map[*Activity]time.Time)
	return activityMap
}

func (tracker *ActivityTracker) registerServices(services ...func(tickerCh chan struct{},
	clickComm chan *Activity) (quit chan struct{})) {

	tracker.activityCh = make(chan *Activity, len(services))    // number based on types of activities being tracked
	tracker.workerTickerCh = make(chan struct{}, len(services)) //this is for all the services, instead of each having their own

	//NO NEED FOR THIS, WE CAN JUST TRY TO RANGE OVER WORKER TICKERCH
	//OR INDIVIDUAL CHANNEL FOR EACH SERVICE
	//ALSO, WHETHER NEED TO PASS VARIABLES TO GO ROUTINES
	for _, service := range services {
		quitCh := service(tracker.workerTickerCh, tracker.activityCh)
		tracker.services = append(tracker.services, quitCh)
	}
}
