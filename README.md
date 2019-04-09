# Activity tracker

[![codecov](https://codecov.io/gh/prashantgupta24/activity-tracker/branch/master/graph/badge.svg)](https://codecov.io/gh/prashantgupta24/activity-tracker) [![Go Report Card](https://goreportcard.com/badge/github.com/prashantgupta24/activity-tracker)](https://goreportcard.com/report/github.com/prashantgupta24/activity-tracker) [![version][version-badge]][RELEASES]

It is a libary that lets you monitor certain activities on your machine, and then sends a [heartbeat](https://github.com/prashantgupta24/activity-tracker#heartbeat-struct) at a periodic (configurable) time detailing all the activity changes during that time. The activities that you want to track are monitored by **pluggable** handlers for those activities and can be added or removed according to your needs. An example of an activity is `MouseCursorActivity`, i.e. whether your mouse cursor was moved or not.

## Installation

` go get -u github.com/prashantgupta24/activity-tracker`

## Usage

```go
heartbeatInterval := 60 //value always in seconds
workerInterval := 5     //seconds

activityTracker := &tracker.Instance{
	HeartbeatInterval: heartbeatInterval,
	WorkerInterval:    workerInterval,
	LogLevel:          logging.Info,
}

//This starts the tracker for all handlers currently implemented. It gives you a channel on
//which you can listen to for heartbeat objects
heartbeatCh := activityTracker.Start()

//if you only want to track certain handlers, you can use StartWithhandlers
//heartbeatCh := activityTracker.StartWithHandlers(handler.MouseClickHandler(), handler.MouseCursorHandler())


select {
case heartbeat := <-heartbeatCh:
	if !heartbeat.WasAnyActivity {
		logger.Infof("no activity detected in the last %v seconds", int(heartbeatInterval))
	} else {
		logger.Infof("activity detected in the last %v seconds.", int(heartbeatInterval))
		logger.Infof("Activity type:\n")
		for activityType, times := range heartbeat.ActivityMap {
			logger.Infof("activityType : %v times: %v\n", activityType, len(times))
		}
	}
}
```

## Output

The above code created a tracker with all (`Mouse-click`, `Mouse-movement`, `screen-change` and `machine-sleep`) handlers activated. The `heartbeat Interval` is set to 60 seconds, i.e. every 60 seconds I received a `heartbeat` which mentioned all activities that were captured.

```sh
INFO[2019-03-30T15:52:01-07:00] starting activity tracker with 60s heartbeat and 5s worker Interval... 

INFO[2019-03-30T15:53:01-07:00] activity detected in the last 60 seconds.    

INFO[2019-03-30T15:53:01-07:00] Activity type:                               
INFO[2019-03-30T15:53:01-07:00] activityType : mouse-click times: 10         
INFO[2019-03-30T15:53:01-07:00] activityType : cursor-move times: 12
INFO[2019-03-30T15:53:01-07:00] activityType : screen-change times: 7
INFO[2019-03-30T15:53:01-07:00] activityType : machine-sleep times: 1
INFO[2019-03-30T15:53:01-07:00] activityType : machine-wake times: 1
```

## How it works

There are 2 primary configs required for the tracker to work:

- `HeartbeatInterval ` 

> The Interval at which you want the heartbeat (in seconds, default 60s)


- `WorkerInterval` 

> The Interval at which you want the checks to happen within a heartbeat (default 60s).

The activity tracker gives you a [heartbeat](https://github.com/prashantgupta24/activity-tracker#heartbeat-struct) object every 60 seconds, that is based on the `HeartbeatInterval`. But there is something else to understand here. In order for the tracker to know how many times an activity occured, like how many times you moved the cursor for example, it needs to query the mouse position every `x` seconds. That's where the `WorkerInterval` comes into play.

The `WorkerInterval` tells the tracker how frequently to check for an activity within a heartbeat. It does that by querying the handler associated with that activity. Let's say you want to know how many times the mouse cursor was moved within 60 seconds. You need to constantly ask the `mouseCursorHandler` every `x` seconds to see if the cursor moved. What you want to do is to start the tracker with the usual 60s `HeartbeatInterval `, configured with a `Mouse-cursor` handler. In this case, you set the `WorkerInterval` to 5 seconds. The tracker will then keep asking the mouse cursor handler every 5 seconds to see if there was a movement, and keep track each time there was a change. At the end of `HeartbeatInterval`, it will construct the `heartbeat` with all the changes and send it.

- If you want to know how many `times` an activity occured within a heartbeat, you might want to set the `WorkerInterval` to a low value, so that it keeps quering the handlers.

- If you are just concerned whether any activity happened within a heartbeat or not, you can set `WorkerInterval` the same as `HeartbeatInterval`. That way, the workers need to check just once before each heartbeat to know if there was any activity registered.


> Note: If the `WorkerInterval` and the `HeartbeatInterval` are set the same, then the `WorkerInterval` always is started a fraction of a second before the `HeartbeatInterval` kicks in. This is done so that when the `heartbeat` is going to be generated at the end of `HeartbeatInterval`, the worker should have done its job of querying each of the handlers before that. 

## Usecase

Suppose you want to track Activities A, B and C on your machine, and you want to know how many times they occured every minute. 

You want a report at the end of every minute saying `Activity A` happened 5 times, `Activity B` happened 3 times and `Activity C` happened 2 times.

First, you need to create a `Handler` for each of those activities. See sections below on how to create one. The main `tracker` object will simply ask each of the handlers every `WorkerInterval` amout of time whether that activity happened or not at that moment.

As another example, let's say you want to monitor whether there was any mouse click on your machine and you want to be notified every 5 minutes. What you do is start the `Activity Tracker` with just the `mouse click` handler and `heartbeat` Interval set to 5 minutes. The `Start` function of the library gives you a channel which receives a `heartbeat` every 5 minutes, and it has details on whether there was a `click` in those 5 minutes, and if yes, the times the click happened.


# Components

### Heartbeat struct

It is the data packet sent from the tracker library to the user.

```go
type Heartbeat struct {
	WasAnyActivity bool
	ActivityMap    map[activity.Type][]time.Time //activity type with its times
	Time           time.Time                     //heartbeat time
}
```

`WasAnyActivity` tells if there was any activity within that time frame
If there was, then the `ActivityMap` will tell you what type of activity it was and what all times it occured.

The `Time` field is the time of the Heartbeat sent (not to be confused with
the activity time, which is the time the activity occured within the `heartbeat`). 

### Tracker

The tracker is the main struct for the library. The fields inside it are:

```go
HeartbeatInterval int //the interval at which you want the heartbeat (in seconds, default 60s)
WorkerInterval    int //the interval at which you want the checks to happen within a heartbeat (in seconds, default 5s)
LogLevel          string //info or debug
LogFormat         string //text or json
```

#### - `HeartbeatInterval` 

The Interval at which you want the heartbeat (in seconds, default 60s)

> The `HeartbeatInterval ` value can be set anywhere between 60 seconds - 300 seconds. Not setting it or setting it to anything other than the allowed range will revert it to default of 60s.

#### - `WorkerInterval` 

The Interval at which you want the checks to happen within a heartbeat (default 60s).

> The `WorkerInterval ` value can be set anywhere between 4 seconds - 60 seconds. It CANNOT be more than `HeartbeatInterval` for obvious reasons. Not setting it or setting it to anything other than the allowed range will revert it to default of 60s.

### State

The `system.State` struct captures the current state of the tracker, and the whole system in general. It is used by some of the handlers to respond to a certain system state. 

It is passed to the handlers when performing the Trigger, so that the handlers can take an informed decision on whether to get activated or not at that instance.

For example, the `sleepHandler` changes the state of the system to sleeping, so that the `mouseCursorHandler` and `mouseClickHandler` don't need to do any work while the system remains in the sleep state.

> Note: It also serves as a way of inter-handler communication.

## Types of handlers

There are 2 types of handlers:

- Push based
- Pull based


The `push` based ones are those that push to the `tracker` object when an activity happened. Examples are the `mouseClickHander` and `machineSleepHandler`. Whenever a mouse-click/machine-sleep happens, it sends the `activity` to the `tracker` object.

The `pull` based ones are those that the `tracker` has to ask the handler to know if there was any activity happening at that moment.
Examples are `mouseCursorHandler` and `screenChangeHandler`. The `asking` is done through the `Trigger` function implemented by handlers.

It is up to you to define how to implement the handler. Some make sense to be pull based, since it is going to be memory intensive to make the mouse cursor movement handler push-based. It made sense to make it `pull` based.

## New pluggable handlers for activities

```go
//Handler interface
Start(*log.Logger, chan *activity.Instance)
Type() activity.Type
Trigger(system.State) //used to activate pull-based handlers
Close()
```
	
Any new type of handler for an activity can be easily added, it just needs to implement the above `Handler` interface and define what `type` of activity it is going to track (also add the new `activity` as well if it's a new activity), that's it! It can be plugged in with the tracker and then the tracker will include those activity checks in its heartbeat.

> Note: Handlers have a many-to-one relationship with activity, i.e. each Handler can be associated with one or more activity (That becomes the value returned by handler's `Type`) On the other hand, each activity should be tracked by only ONE handler (which makes sense).
> As a fail-safe, if the tracker is started with more than one handler tracking the same activity, then only 1 handler will get registered for that activity.

## Currently supported list of activities/handlers


### Activities

```go
MouseCursorMovement Type = "cursor-move"
MouseClick          Type = "mouse-click"
MachineSleep        Type = "machine-sleep"
MachineWake         Type = "machine-wake"
```

### Corresponding handlers

```go
mouseCursorHandler
mouseClickHandler
screenChangeHandler
machineSleepHandler
```
	
- Mouse click (whether any mouse click happened during the time frame)
- Mouse cursor movement (whether the mouse cursor was moved during the time frame)
- Machine sleep/wake handler (**this is added by default for fail-safe measures**)
- Screen change handler (whether the active window was changed)
## Example

Check out the example [here](https://github.com/prashantgupta24/activity-tracker/blob/master/example/example.go)

[version-badge]: https://img.shields.io/github/release/prashantgupta24/activity-tracker.svg
[RELEASES]: https://github.com/prashantgupta24/activity-tracker/releases
