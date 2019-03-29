package activity

//Type of activity as defined below
type Type string

/*
MouseCursorMovement
MouseClick
ScreenChange
TestActivity

These are the types of activities the tracker currently supports
*/
const (
	MouseCursorMovement Type = "cursor-move"
	MouseClick          Type = "mouse-click"
	ScreenChange        Type = "screen-change"
	TestActivity        Type = "test-activity"
)

//Instance is an instance of Activity
type Instance struct {
	Type Type
}
