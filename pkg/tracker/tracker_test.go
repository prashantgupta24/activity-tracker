package tracker

import (
	"sync"
	"testing"
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/handler"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"

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
}

//Run once before each test
func (suite *TestTracker) SetupTest() {
	frequency := 1 //second

	suite.tracker = &Instance{
		Frequency: frequency,
		isTest:    true,
	}
}

func (suite *TestTracker) TestTrackerValidateFrequency() {
	t := suite.T()

	var tracker *Instance

	//testing with empty
	tracker = &Instance{}
	frequency, _ := tracker.validateFrequency()
	assert.Equal(t, time.Duration(defaultFrequency), frequency, "tracker should get default frequency")

	//testing with 0
	tracker = &Instance{
		Frequency: 0,
	}
	frequency, _ = tracker.validateFrequency()
	assert.Equal(t, time.Duration(defaultFrequency), frequency, "tracker should get default frequency")

	//testing with -1
	tracker = &Instance{
		Frequency: -1,
	}
	frequency, _ = tracker.validateFrequency()
	assert.Equal(t, time.Duration(defaultFrequency), frequency, "tracker should get default frequency")

	//testing with min
	tracker = &Instance{
		Frequency: minValFrequency,
	}
	frequency, _ = tracker.validateFrequency()
	assert.Equal(t, time.Duration(minValFrequency), frequency, "tracker should retain original frequency")

	//testing with max
	tracker = &Instance{
		Frequency: maxValFrequency,
	}
	frequency, _ = tracker.validateFrequency()
	assert.Equal(t, time.Duration(maxValFrequency), frequency, "tracker should retain original frequency")

	//testing with test instance = false
	tracker = &Instance{
		Frequency: 1,
		isTest:    false,
	}
	frequency, _ = tracker.validateFrequency()
	assert.Equal(t, time.Duration(defaultFrequency), frequency, "tracker should get default frequency")

	//testing with test instance = true
	tracker = &Instance{
		Frequency: 1,
		isTest:    true,
	}
	frequency, _ = tracker.validateFrequency()
	assert.Equal(t, time.Duration(1), frequency, "tracker should retain original frequency since it is a test")
}
func (suite *TestTracker) TestDupHandlerRegistration() {
	t := suite.T()
	tracker := suite.tracker

	tracker.StartWithHandlers(handler.TestHandler(),
		handler.TestHandler())

	assert.Equal(t, 1, len(tracker.handlers), "duplicate handlers should not be registered")
}

func (suite *TestTracker) TestActivityHandlersNumEqual() {
	t := suite.T()
	numActivities := len(suite.activities)
	numHandlers := len(getAllHandlers())

	assert.Equal(t, numHandlers, numActivities, "tracker should have equal handlers and activities")
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

	numHandlers := len(getAllHandlers())
	heartbeatCh := tracker.Start()

	assert.Equal(t, numHandlers, len(tracker.handlers), "tracker should have started with %v handlers by default", numHandlers)

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
