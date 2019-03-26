package tracker

import (
	"sync"
	"testing"
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/service"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestTracker struct {
	suite.Suite
	activities []*activity.Type
	tracker    *Instance
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TestTracker))
}

//Run once before all tests
func (suite *TestTracker) SetupSuite() {

	suite.activities = append(suite.activities, &activity.Type{
		ActivityType: activity.MouseClick,
	})
	suite.activities = append(suite.activities, &activity.Type{
		ActivityType: activity.MouseCursorMovement,
	})
	suite.activities = append(suite.activities, &activity.Type{
		ActivityType: activity.ScreenChange,
	})
}

//Run once before each test
func (suite *TestTracker) SetupTest() {
	frequency := 1

	suite.tracker = &Instance{
		Frequency: frequency,
	}
}

func (suite *TestTracker) TestDupServiceRegistration() {
	t := suite.T()
	tracker := suite.tracker

	tracker.StartWithServices(service.TestHandler(),
		service.TestHandler())

	assert.Equal(t, 1, len(tracker.services), "duplicate services should not be registered")
}

func (suite *TestTracker) TestActivityServicesNumEqual() {
	t := suite.T()
	numActivities := len(suite.activities)
	numServices := len(getAllServiceHandlers())

	assert.Equal(t, numServices, numActivities, "tracker should have equal services and activities")
}

func (suite *TestTracker) TestActivitiesOneByOne() {
	t := suite.T()
	tracker := suite.tracker

	heartbeatCh := tracker.StartWithServices()

	//send one activity at a time, then wait for heartbeat to acknowledge it
	for _, sentActivity := range suite.activities {
		//fmt.Printf("sending %v activity to tracker\n", sentActivity.ActivityType)
		tracker.activityCh <- sentActivity
		select {
		case heartbeat := <-heartbeatCh:
			assert.NotNil(t, heartbeat)
			assert.True(t, heartbeat.WasAnyActivity, "there should have been activity")

			for activity, time := range heartbeat.Activity {
				assert.Equal(t, sentActivity, activity, "should be of %v activity type", sentActivity)
				assert.NotNil(t, time, "time of activity should not be nil")
			}
		}
	}
}

func (suite *TestTracker) TestActivitiesAllAtOnce() {
	t := suite.T()
	tracker := suite.tracker

	heartbeatCh := tracker.StartWithServices()

	//send all activities
	for _, sentActivity := range suite.activities {
		tracker.activityCh <- sentActivity
	}
	select {
	case heartbeat := <-heartbeatCh:
		assert.NotNil(t, heartbeat)
		assert.True(t, heartbeat.WasAnyActivity, "there should have been activity")
		assert.Equal(t, len(suite.activities), len(heartbeat.Activity), "tracker should registered %v activities ", len(suite.activities))
		for activity, time := range heartbeat.Activity {
			assert.Contains(t, suite.activities, activity, "should contain %v activity type", activity.ActivityType)
			assert.NotNil(t, time, "time of activity should not be nil")
		}
	}
}

func (suite *TestTracker) TestServiceTestHandler() {
	t := suite.T()
	tracker := suite.tracker

	heartbeatCh := tracker.StartWithServices(service.TestHandler())

	select {
	case heartbeat := <-heartbeatCh:
		assert.NotNil(t, heartbeat)
		assert.True(t, heartbeat.WasAnyActivity, "there should have been activity")
		testActivity := &activity.Type{
			ActivityType: activity.TestActivity,
		}
		for activity, time := range heartbeat.Activity {
			assert.Equal(t, activity, testActivity, "should be of test activity type")
			assert.NotNil(t, time, "time of activity should not be nil")
		}
	}
}

func (suite *TestTracker) TestTrackerStartAndQuit() {
	t := suite.T()
	tracker := suite.tracker

	numServices := len(getAllServiceHandlers())
	heartbeatCh := tracker.Start()

	assert.Equal(t, numServices, len(tracker.services), "tracker should have started with %v services by default", numServices)

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
