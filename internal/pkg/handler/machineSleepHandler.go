package handler

import (
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
	"github.com/prashantgupta24/activity-tracker/pkg/system"
	"github.com/prashantgupta24/mac-sleep-notifier/notifier"
	log "github.com/sirupsen/logrus"
)

const (
	machineSleep = activity.MachineSleep
	machineWake  = activity.MachineWake
)

//MachineSleepHanderStruct is a handler for machine sleep/awake related events
type MachineSleepHanderStruct struct {
	sleepHandlerLogger *log.Entry
	quit               chan struct{}
}

//Start the handler
func (m *MachineSleepHanderStruct) Start(logger *log.Logger, activityCh chan *activity.Instance) {
	m.quit = make(chan struct{})
	m.sleepHandlerLogger = logger.WithFields(log.Fields{
		"method": "machine-sleep-handler",
	})

	go func() {
		notifierCh := notifier.GetInstance().Start()
		for {
			select {
			case notification := <-notifierCh:
				if notification.Type == notifier.Awake {
					m.sleepHandlerLogger.Debug("machine awake")
					activityCh <- &activity.Instance{
						Type: machineWake,
						State: &system.State{
							IsSystemSleep: false,
						},
					}
				} else {
					if notification.Type == notifier.Sleep {
						m.sleepHandlerLogger.Debug("machine sleeping")
						activityCh <- &activity.Instance{
							Type: machineSleep,
							State: &system.State{
								IsSystemSleep: true,
							},
						}
					}
				}
			case <-m.quit:
				logger.Infof("stopping sleep handler")
				return
			}
		}
	}()
}

//MachineSleepHandler returns an instance of the struct
func MachineSleepHandler() *MachineSleepHanderStruct {
	return &MachineSleepHanderStruct{}
}

//Type returns the type of handler
func (m *MachineSleepHanderStruct) Type() activity.Type {
	return machineSleep
}

//Trigger the handler - empty since it's a push based handler
func (m *MachineSleepHanderStruct) Trigger(system.State) {}

//Close closes the handler
func (m *MachineSleepHanderStruct) Close() {
	m.quit <- struct{}{}
}
