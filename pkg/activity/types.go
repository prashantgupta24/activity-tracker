package activity

type activityType string

const (
	MOUSE_CURSOR_MOVEMENT activityType = "cursor-move"
	MOUSE_CLICK           activityType = "mouse-click"
	SCREEN_CHANGE         activityType = "screen-change"
	TEST_ACTIVITY         activityType = "test-activity"
)

type Type struct {
	ActivityType activityType
}
