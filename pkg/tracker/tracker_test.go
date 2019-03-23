package tracker

import (
	"testing"

	"github.com/prashantgupta24/activity-tracker/pkg/activity"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestTracker struct {
	suite.Suite
	tracker *Instance
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TestTracker))
}

//Run once before all tests
func (suite *TestTracker) SetupSuite() {

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

	_ = tracker.StartWithServices(service.TestHandler(),
		service.TestHandler())

	assert.Equal(t, 1, len(tracker.services), "duplicate services should not be registered")
}

func (suite *TestTracker) TestServiceHandler() {
	t := suite.T()
	tracker := suite.tracker

	heartbeatCh := tracker.StartWithServices(service.TestHandler())

	select {
	case heartbeat := <-heartbeatCh:
		assert.NotNil(t, heartbeat)
		assert.True(t, heartbeat.IsActivity, "should be an activity")
		testActivity := &activity.Type{
			ActivityType: activity.TEST_ACTIVITY,
		}
		for activity, time := range heartbeat.Activity {
			assert.Equal(t, activity, testActivity, "should be of test activity type")
			assert.NotNil(t, time, "time of activity should not be nil")

		}
		// fmt.Println("heartbeat time : ", heartbeat.Time)
		// fmt.Println("activity time : ", heartbeat.Activity)
	}
}
