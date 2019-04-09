package system

/*
State captures the current state of the tracker, and the whole system in general.

It is passed to the handlers when performing the Trigger, so that the handlers
can take an informed decision on whether to get activated or not at that instance.

For example: The mouseCursorHandler does not need to do anything if the state of
the system is in sleep state.

It can serve as a way of inter-handler communication
*/
type State struct {
	IsSystemSleep bool //this state variable captures whether the system is sleeping or not
}
