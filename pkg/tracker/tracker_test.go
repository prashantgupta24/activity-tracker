package tracker

import (
	"testing"

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
		ActivityType: activity.MOUSE_CLICK,
	})
	suite.activities = append(suite.activities, &activity.Type{
		ActivityType: activity.MOUSE_CURSOR_MOVEMENT,
	})
	suite.activities = append(suite.activities, &activity.Type{
		ActivityType: activity.SCREEN_CHANGE,
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

	tracker.StartWithServices(service.MouseClickHandler(),
		service.MouseClickHandler())

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

func (suite *TestTracker) TestAllServicesStart() {
	//t := suite.T()

	// activityCh := make(chan *activity.Type)
	// logger := logging.New()
	// for _, service := range getAllServiceHandlers() {
	// 	service.Start(logger, activityCh)
	// }
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
			ActivityType: activity.TEST_ACTIVITY,
		}
		for activity, time := range heartbeat.Activity {
			assert.Equal(t, activity, testActivity, "should be of test activity type")
			assert.NotNil(t, time, "time of activity should not be nil")
		}
	}
}

// func (suite *TestTracker) TestServiceClose() {
// 	//t := suite.T()
// 	tracker := suite.tracker

// 	testHandler := service.TestHandler()
// 	tracker.StartWithServices(testHandler)
// 	tracker.Quit()
// }
