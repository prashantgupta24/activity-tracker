# Activity tracker

[![codecov](https://codecov.io/gh/prashantgupta24/activity-tracker/branch/master/graph/badge.svg)](https://codecov.io/gh/prashantgupta24/activity-tracker) [![Go Report Card](https://goreportcard.com/badge/github.com/prashantgupta24/activity-tracker)](https://goreportcard.com/report/github.com/prashantgupta24/activity-tracker) [![version][version-badge]][RELEASES]

It is a libary that lets you monitor certain activities on your machine, and sends a heartbeat at a periodic (configurable) time detailing all the activity changes during that time. The activities that you want to monitor are **pluggable** handlers for those activities and can be added or removed according to your needs.

## Installation

` go get -u github.com/prashantgupta24/activity-tracker`

## Usage


	heartbeatFrequency := 60 //value always in seconds
	workerFrequency := 5     //seconds

	activityTracker := &tracker.Instance{
		HeartbeatFrequency: heartbeatFrequency,
		WorkerFrequency:    workerFrequency,
		LogLevel:           logging.Debug,
	}

	//This starts the tracker for all handlers
	heartbeatCh := activityTracker.Start()

	//if you only want to track certain handlers, you can use StartWithhandlers
	//heartbeatCh := activityTracker.StartWithHanders(handler.MouseClickHandler(), handler.MouseCursorHandler())


		select {
		case heartbeat := <-heartbeatCh:
			if !heartbeat.WasAnyActivity {
				logger.Infof("no activity detected in the last %v seconds\n\n\n", int(heartbeatFrequency))
			} else {
				logger.Infof("activity detected in the last %v seconds.", int(heartbeatFrequency))
				logger.Infof("Activity type:\n")
				for activityType, times := range heartbeat.ActivityMap {
					logger.Infof("activityType : %v times: %v\n", activityType, len(times))
				}
			}
		}

## Output

The above code created a tracker with all (`Mouse-click`, `Mouse-movement` and `Screen-Change`) handlers activated. The `heartbeat frequency` is set to 60 seconds, i.e. every 60 seconds I received a `heartbeat` which mentioned all activities that were captured.

```
INFO[2019-03-30T15:52:01-07:00] starting activity tracker with 60s heartbeat and 5s worker frequency... 

INFO[2019-03-30T15:53:01-07:00] activity detected in the last 60 seconds.    

INFO[2019-03-30T15:53:01-07:00] Activity type:                               
INFO[2019-03-30T15:53:01-07:00] activityType : screen-change times: 7        
INFO[2019-03-30T15:53:01-07:00] activityType : mouse-click times: 10         
INFO[2019-03-30T15:53:01-07:00] activityType : cursor-move times: 12   
```

## How it works

There are 2 primary configs required for the tracker to work:

- `HeartbeatFrequency ` 

> The frequency at which you want the heartbeat (in seconds, default 60s)


- `WorkerFrequency` 

> The frequency at which you want the checks to happen within a heartbeat (default 60s).


The activity tracker gives you a `heartbeat` object every 60 seconds, that is based on the `HeartbeatFrequency `. but there is something else to understand here. In order for the tracker to know how many times an activity occured, or how many times you moved the cursor for example, it needs to query the mouse movement library `n` number of times. That's where the `WorkerFrequency` comes into play.

The `WorkerFrequency` tells the tracker how many times to query for each of the handlers in the tracker within a heartbeat. If you are just concerned whether any activity happened within a heartbeat or not, you can set it to the same as `HeartbeatFrequency`. 

If you want to know how many times an activity occured within a heartbeat, you might want to set the `WorkerFrequency` to a low value, so that it keeps quering the handlers.


## Another example

Suppose you want to track Activities A, B and C on your machine, and you want a heartbeat every 5 minutes. What it would do then is to send you heartbeats every 5 minutes, and each heartbeat would contain whether any of A, B or C occured within those 5 minutes, and if so, at what times.

As another example, let's say you want to monitor whether there was any mouse click on your machine and you want to be monitor every 5 minutes. What you do is start the `Activity Tracker` with just the `mouse click` handler and `heartbeat` frequency set to 5 minutes. The `Start` function of the library gives you a channel which receives a `heartbeat` every 5 minutes, and it has details on whether there was a `click` in those 5 minutes, and if yes, the times the click happened.



# Components

### Heartbeat struct

It is the data packet sent from the tracker library to the user.

	type Heartbeat struct {
		WasAnyActivity bool //whether any activity was detected 		
		ActivityMap       map[activity.Type][]time.Time //activity type with its times
		Time           time.Time                    //heartbeat time
	}

`WasAnyActivity` tells if there was any activity within that time frame
If there was, then the `ActivityMap` will tell you what type of activity it was and what all times it occured.

The `Time` field is the time of the Heartbeat sent (not to be confused with
the activity time, which is the time the activity occured within the time frame)

### Tracker

The tracker is the main struct for the library. 

	//Instance is an instance of the tracker
	HeartbeatFrequency int //the frequency at which you want the heartbeat (in seconds, default 60s)
	WorkerFrequency    int //therequency at which you want the checks to happen within a heartbeat (in seconds, default 5s)
	LogLevel           string
	LogFormat          string


#### - `HeartbeatFrequency ` 

The frequency at which you want the heartbeat (in seconds, default 60s)

##### Note: The `HeartbeatFrequency ` value can be set anywhere between 60 seconds - 300 seconds. Not setting it or setting it to anything other than the allowed range will revert it to default of 60s.

#### - `WorkerFrequency` 

The frequency at which you want the checks to happen within a heartbeat (default 60s).

##### Note: The `WorkerFrequency ` value can be set anywhere between 4 seconds - 60 seconds. It CANNOT be more than `HeartbeatFrequency` for obvious reasons. Not setting it or setting it to anything other than the allowed range will revert it to default of 60s.


## Relationship between Activity and Handler
	
Activity and Handler have a 1-1 mapping, i.e. each handler can only handle one type of activity, and vice-versa, each activity should be handled by one handler only.

The `Type` in the `Handler` interface determines the type of activity the particular handler handles.

## New pluggable handlers for activities


	//Instance is the main interface for a Handler for the tracker
	type Instance interface {
		Start(*log.Logger, chan *activity.Instance)
		Type() activity.Type
		Trigger()
		Close()
	}
	
Any new type of handler for an activity can be easily added, it just needs to implement the above `Handler` interface and define what type of activity it is going to track, that's it! It can be plugged in with the tracker and then the tracker will include those activity checks in its heartbeat.


## Currently supported list of activities/handlers


### Activities

	MouseCursorMovement Type = "cursor-move"
	MouseClick          Type = "mouse-click"
	ScreenChange        Type = "screen-change"

### Corresponding handlers

	mouseCursorHandler
	mouseClickHandler
	screenChangeHandler
	
	
- Mouse click (whether any mouse click happened during the time frame)
- Mouse cursor movement (whether the mouse cursor was moved during the time frame)
- Screen change (whether the screen was changed anytime within that time frame)



[version-badge]: https://img.shields.io/github/release/prashantgupta24/activity-tracker.svg
[RELEASES]: https://github.com/prashantgupta24/activity-tracker/releases
