# Activity tracker

[![codecov](https://codecov.io/gh/prashantgupta24/activity-tracker/branch/master/graph/badge.svg)](https://codecov.io/gh/prashantgupta24/activity-tracker) [![Go Report Card](https://goreportcard.com/badge/github.com/prashantgupta24/activity-tracker)](https://goreportcard.com/report/github.com/prashantgupta24/activity-tracker) [![version][version-badge]][RELEASES]

It is a libary that lets you monitor certain activities on your machine, and sends a heartbeat at a periodic (configurable) time detailing all the activity changes during that time. The activities that you want to monitor are **pluggable** handlers for those activities and can be added or removed according to your needs.

## Installation

` go get -u github.com/prashantgupta24/activity-tracker`
## Usage


	frequency := 60 //value always in seconds

	activityTracker := &tracker.Instance{
		Frequency: frequency,
		LogLevel:  logging.Info,
	}

	//This starts the tracker for all handlers
	heartbeatCh := activityTracker.Start()

	//if you only want to track certain handlers, you can use StartWithhandlers
	//heartbeatCh := activityTracker.StartWithHanders(handler.MouseClickHandler(), handler.MouseCursorHandler())


		select {
		case heartbeat := <-heartbeatCh:
			if !heartbeat.WasAnyActivity {
				logger.Infof("no activity detected in the last %v seconds\n\n\n", int(frequency))
			} else {
				logger.Infof("activity detected in the last %v seconds.", int(frequency))
				logger.Infof("Activity type:\n")
				for activityType, times := range heartbeat.ActivityMap {
					logger.Infof("activityType : %v times: %v\n", activityType, len(times))
				}
			}
		}

## Output

The above code created a tracker with all (`Mouse-click`, `Mouse-movement` and `Screen-Change`) handlers activated. The `heartbeat` frequency is set to 60 seconds, i.e. every 60 seconds I received a `heartbeat` which mentioned all activities that were captured.

```
INFO[2019-03-30T14:50:17-07:00] starting activity tracker with 60 second frequency ... 

INFO[2019-03-30T14:51:17-07:00] activity detected in the last 60 seconds.    

INFO[2019-03-30T14:51:17-07:00] Activity type:                               
INFO[2019-03-30T14:51:17-07:00] activityType : mouse-click times: 3          
INFO[2019-03-30T14:51:17-07:00] activityType : cursor-move times: 5          
INFO[2019-03-30T14:51:17-07:00] activityType : screen-change times: 3  
```

## How it works

The activity tracker gives you a `heartbeat` object every 60 seconds, but there is something else to understand here. In order for the tracker to know how many times you moved the cursor, it needs to query the mouse movement library which gets the actual mouse cursor position. Ideally, I should be checking this every second to see if the mouse moved, but that will prove very 

## Example

Suppose you want to track Activities A, B and C on your machine, and you want the tracker to monitor every 5 minutes. What it would do then is to send you heartbeats every 5 minutes, and each heartbeat would contain whether any of A, B or C occured within those 5 minutes, and if so, at what times.

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
If there was, then the ActivityMap will tell you what type of activity
it was and what times it occured.

The `Time` field is the time of the Heartbeat sent (not to be confused with
the activity time, which is the time the activity occured within the time frame)

### Tracker instance

	//Instance is an instance of the tracker
	type Instance struct {
		Frequency  int //the frequency at which you want the heartbeat (in seconds)
		LogLevel   string
		LogFormat  string

The tracker is the main struct for the library. The `Frequency` is the main component which can be changed according to need.

##### Note: The `frequency` value can be set anywhere between 60 seconds - 300 seconds. Not setting it or setting it to anything other than the allowed range will revert it to default of 60s.

### Activity and Handler

	//Instance is an instance of Activity
	type Instance struct {
		Type Type
	}
	
	//Instance is the main interface for a Handler for the tracker
	type Instance interface {
		Start(*log.Logger, chan *activity.Instance)
		Type() activity.Type
		Trigger()
		Close()
	}
	
Activity and Handler have a 1-1 mapping, i.e. each handler can only handle one type of activity.

The `Type` in the `Handler` interface determines the type of activity the particular handler handles.

## Currently supported list of handlers

- Mouse click (whether any mouse click happened during the time frame)
- Mouse cursor movement (whether the mouse cursor was moved during the time frame)
- Screen change (whether the screen was changed anytime within that time frame)

## New pluggable handlers for activities

Any new type of handler for an activity can be easily added, it just needs to implement the `Handler` interface and define what type of activity it is going to track, that's it! It can be plugged in with a tracker and then the tracker will include those activity checks in its heartbeat.


[version-badge]: https://img.shields.io/github/release/prashantgupta24/activity-tracker.svg
[RELEASES]: https://github.com/prashantgupta24/activity-tracker/releases
