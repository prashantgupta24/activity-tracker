package activity

type activityType string

/*
MouseCursorMovement
MouseClick
ScreenChange
TestActivity

These are the types of activities the tracker currently supports
*/
const (
	MouseCursorMovement activityType = "cursor-move"
	MouseClick          activityType = "mouse-click"
	ScreenChange        activityType = "screen-change"
	TestActivity        activityType = "test-activity"
)

//Type gets the type of Activity
type Type struct {
	ActivityType activityType
}
