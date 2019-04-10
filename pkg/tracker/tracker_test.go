package tracker

import (
	"sync"
	"testing"
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/handler"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
	"github.com/prashantgupta24/activity-tracker/pkg/system"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestTracker struct {
	suite.Suite
	activities []activity.Type
	tracker    *Instance
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TestTracker))
}

//Run once before all tests
func (suite *TestTracker) SetupSuite() {
	suite.activities = append(suite.activities, activity.MouseClick)
	suite.activities = append(suite.activities, activity.MouseCursorMovement)
	suite.activities = append(suite.activities, activity.ScreenChange)
	suite.activities = append(suite.activities, activity.MachineWake)
	suite.activities = append(suite.activities, activity.MachineSleep)
}

//Run once before each test
func (suite *TestTracker) SetupTest() {
	heartbeatInterval := 1 //second

	suite.tracker = &Instance{
		HeartbeatInterval: heartbeatInterval,
		isTest:            true,
		//LogLevel:           "debug",
	}
}

func (suite *TestTracker) TestTrackerValidateInterval() {
	t := suite.T()

	var tracker *Instance

	//testing with empty
	tracker = &Instance{}
	heartbeatInterval, workerInterval := tracker.validateIntervals()
	assert.Equal(t, time.Duration(defaultHInterval), heartbeatInterval, "tracker should get default Interval")
	assert.Equal(t, time.Duration(defaultWInterval), workerInterval, "tracker should get default Interval")

	//testing with 0
	tracker = &Instance{
		HeartbeatInterval: 0,
		WorkerInterval:    0,
	}
	heartbeatInterval, workerInterval = tracker.validateIntervals()
	assert.Equal(t, time.Duration(defaultHInterval), heartbeatInterval, "tracker should get default Interval")
	assert.Equal(t, time.Duration(defaultWInterval), workerInterval, "tracker should get default Interval")

	//testing with -1
	tracker = &Instance{
		HeartbeatInterval: -1,
		WorkerInterval:    -1,
	}
	heartbeatInterval, workerInterval = tracker.validateIntervals()
	assert.Equal(t, time.Duration(defaultHInterval), heartbeatInterval, "tracker should get default Interval")
	assert.Equal(t, time.Duration(defaultWInterval), workerInterval, "tracker should get default Interval")

	//testing with min
	tracker = &Instance{
		HeartbeatInterval: minHInterval,
		WorkerInterval:    minWInterval,
	}
	heartbeatInterval, workerInterval = tracker.validateIntervals()
	assert.Equal(t, time.Duration(minHInterval), heartbeatInterval, "tracker should retain original Interval")
	assert.Equal(t, time.Duration(minWInterval), workerInterval, "tracker should retain original Interval")

	//testing with max
	tracker = &Instance{
		HeartbeatInterval: maxHInterval,
		WorkerInterval:    maxWInterval,
	}
	heartbeatInterval, workerInterval = tracker.validateIntervals()
	assert.Equal(t, time.Duration(maxHInterval), heartbeatInterval, "tracker should retain original Interval")
	assert.Equal(t, time.Duration(maxWInterval), workerInterval, "tracker should retain original Interval")

	//testing with test instance = false
	tracker = &Instance{
		HeartbeatInterval: 1,
		isTest:            false,
	}
	heartbeatInterval, workerInterval = tracker.validateIntervals()
	assert.Equal(t, time.Duration(defaultHInterval), heartbeatInterval, "tracker should get default Interval")
	assert.Equal(t, time.Duration(defaultWInterval), workerInterval, "tracker should get default Interval")

	//testing with test instance = true
	tracker = &Instance{
		HeartbeatInterval: 1,
		isTest:            true,
	}
	heartbeatInterval, workerInterval = tracker.validateIntervals()
	assert.Equal(t, time.Duration(1), heartbeatInterval, "tracker should retain original Interval since it is a test")
	assert.Equal(t, heartbeatInterval, workerInterval, "worker should match heartbeat since it is a test")
}
func (suite *TestTracker) TestDupHandlerRegistration() {
	t := suite.T()
	tracker := suite.tracker

	tracker.StartWithHandlers(handler.TestHandler(),
		handler.TestHandler())

	assert.Equal(t, 1, len(tracker.handlers), "duplicate handlers should not be registered")
}

func (suite *TestTracker) TestActivitiesOneByOne() {
	t := suite.T()
	tracker := suite.tracker

	heartbeatCh := tracker.StartWithHandlers()

	//send one activity at a time, then wait for heartbeat to acknowledge it
	for _, sentActivityType := range suite.activities {
		tracker.activityCh <- &activity.Instance{
			Type: sentActivityType,
		}
		select {
		case heartbeat := <-heartbeatCh:
			assert.NotNil(t, heartbeat)
			assert.True(t, heartbeat.WasAnyActivity, "there should have been activity")

			for activityType, time := range heartbeat.ActivityMap {
				assert.Equal(t, sentActivityType, activityType, "should be of %v activity type", sentActivityType)
				assert.NotNil(t, time, "time of activity should not be nil")
			}
		}
	}
}

func (suite *TestTracker) TestActivitiesAllAtOnce() {
	t := suite.T()
	tracker := suite.tracker

	heartbeatCh := tracker.StartWithHandlers()

	//send all activities
	for _, sentActivityType := range suite.activities {
		tracker.activityCh <- &activity.Instance{
			Type: sentActivityType,
		}
	}
	select {
	case heartbeat := <-heartbeatCh:
		assert.NotNil(t, heartbeat)
		assert.True(t, heartbeat.WasAnyActivity, "there should have been activity")
		assert.Equal(t, len(suite.activities), len(heartbeat.ActivityMap), "tracker should registered %v activities ", len(suite.activities))
		for activityType, time := range heartbeat.ActivityMap {
			assert.Contains(t, suite.activities, activityType, "should contain %v activity type", activityType)
			assert.NotNil(t, time, "time of activity should not be nil")
		}
	}
}

func (suite *TestTracker) TestMultipleActivities() {
	t := suite.T()
	tracker := suite.tracker

	heartbeatCh := tracker.StartWithHandlers()

	timesToSend := 2
	//send all activities timesToSend times
	for i := 0; i < timesToSend; i++ {
		for _, sentActivityType := range suite.activities {
			tracker.activityCh <- &activity.Instance{
				Type: sentActivityType,
			}
		}
	}

	select {
	case heartbeat := <-heartbeatCh:
		assert.NotNil(t, heartbeat)
		assert.True(t, heartbeat.WasAnyActivity, "there should have been activity")
		assert.Equal(t, len(suite.activities), len(heartbeat.ActivityMap), "tracker should have registered %v activities ", len(suite.activities))
		for activityType, time := range heartbeat.ActivityMap {
			//fmt.Printf("activityType : %v times: %v\n", activityType, time)
			assert.Contains(t, suite.activities, activityType, "should contain %v activityType", activityType)
			assert.NotNil(t, time, "time of activity should not be nil")
			assert.Equal(t, timesToSend, len(time), "activity should be recorded %v times", timesToSend)
		}
	}
}
func (suite *TestTracker) TestTestHandler() {
	t := suite.T()
	tracker := suite.tracker

	heartbeatCh := tracker.StartWithHandlers(handler.TestHandler())

	select {
	case heartbeat := <-heartbeatCh:
		assert.NotNil(t, heartbeat)
		assert.True(t, heartbeat.WasAnyActivity, "there should have been activity")
		for activityType, time := range heartbeat.ActivityMap {
			assert.Equal(t, activity.TestActivity, activityType, "should be of test activity type")
			assert.NotNil(t, time, "time of activity should not be nil")
		}
	}
}

func (suite *TestTracker) TestTrackerStartAndQuit() {
	t := suite.T()
	tracker := suite.tracker

	heartbeatCh := tracker.StartWithHandlers(handler.TestHandler())
	//close
	var wg sync.WaitGroup
	isHeartbeatStopped := make(chan bool)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		for heartbeat := range heartbeatCh {
			assert.IsType(t, &Heartbeat{}, heartbeat, "type not equal")
		}
		wg.Done()
	}(&wg)

	tracker.Quit()

	go func(isHeartbeatStopped chan bool) {
		wg.Wait()
		isHeartbeatStopped <- true
	}(isHeartbeatStopped)

	select {
	case <-time.After(time.Second):
		assert.Fail(t, "heartbeat should have stopped after quit")
	case val := <-isHeartbeatStopped:
		assert.True(t, val)
	}
}

func (suite *TestTracker) TestValidateHandlers() {
	t := suite.T()
	machineSleepHandler := handler.MachineSleepHandler()

	//case 1
	handlers := []handler.Instance{
		handler.MouseClickHandler(), handler.MouseCursorHandler(),
		handler.ScreenChangeHandler(),
	}
	validatedHandlers := validateHandlers(handlers...)
	assert.Contains(t, validatedHandlers, machineSleepHandler, "validateHandler() should add machine sleep handler by default")

	//case 2
	handlers = []handler.Instance{
		handler.MachineSleepHandler(),
	}
	validatedHandlers = validateHandlers(handlers...)
	assert.Contains(t, validatedHandlers, machineSleepHandler, "validateHandler() should add machine sleep handler by default")
}

func (suite *TestTracker) TestValidateHandlersOnStart() {
	t := suite.T()
	machineSleepHandlerType := handler.MachineSleepHandler().Type()

	tracker := suite.tracker

	tracker.StartWithHandlers(handler.TestHandler())
	assert.NotContains(t, tracker.handlers, machineSleepHandlerType, "validateHandler() should not add machine sleep handler type in test")
}

func (suite *TestTracker) TestTrackerStateChange() {
	t := suite.T()
	tracker := suite.tracker
	tracker.StartWithHandlers()

	oldState := tracker.getTrackerSystemState()
	assert.False(t, oldState.IsSystemSleep, "system sleep state should be false by default")

	tracker.activityCh <- &activity.Instance{
		Type: activity.TestActivity,
		State: &system.State{
			IsSystemSleep: true,
		},
	}

	newState := tracker.getTrackerSystemState()
	assert.True(t, newState.IsSystemSleep, "system sleep state should be over-written to true")

}

func (suite *TestTracker) TestTrackerStateChangeByValue() {
	t := suite.T()
	tracker := suite.tracker
	tracker.StartWithHandlers()

	oldState := tracker.getTrackerSystemState()
	oldState.IsSystemSleep = true //this should not affect the state object since it is a copy

	newState := tracker.getTrackerSystemState()
	assert.False(t, newState.IsSystemSleep, "system sleep state should not change")
}
