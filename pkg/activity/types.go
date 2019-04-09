package activity

import "github.com/prashantgupta24/activity-tracker/pkg/system"

//Type of activity as defined below
type Type string

/*
MouseCursorMovement
MouseClick
ScreenChange
TestActivity
MachineSleep
MachineWake

These are the types of activities the tracker currently supports
*/
const (
	MouseCursorMovement Type = "cursor-move"
	MouseClick          Type = "mouse-click"
	ScreenChange        Type = "screen-change"
	TestActivity        Type = "test-activity"
	MachineSleep        Type = "machine-sleep"
	MachineWake         Type = "machine-wake"
)

//Instance is an instance of Activity
type Instance struct {
	Type  Type
	State *system.State //optional, only needed when a handler needs to change the state of the system
}
