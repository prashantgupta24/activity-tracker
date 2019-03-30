# Activity tracker

[![codecov](https://codecov.io/gh/prashantgupta24/activity-tracker/branch/master/graph/badge.svg)](https://codecov.io/gh/prashantgupta24/activity-tracker) [![Go Report Card](https://goreportcard.com/badge/github.com/prashantgupta24/activity-tracker)](https://goreportcard.com/report/github.com/prashantgupta24/activity-tracker) [![version][version-badge]][RELEASES]

It is a libary that lets you monitor certain activities on your machine, and sends a heartbeat at a periodic (configurable) time detailing all the activity changes during that time. The activities that you want to monitor are **pluggable** handlers for those activities and can be added or removed according to your needs.

## Installation

## Usage

## Example

Suppose you want to track Activities A, B and C on your machine, and you want the tracker to monitor every 5 minutes. What it would do then is to send you heartbeats every 5 minutes, and each heartbeat would contain whether any of A, B or C occured within those 5 minutes, and if so, at what times.

As another example, let's say you want to monitor whether there was any mouse click on your machine and you want to be monitor every 5 minutes. What you do is start the `Activity Tracker` with just the `mouse click` handler and `heartbeat` frequency set to 5 minutes. The `Start` function of the library gives you a channel which receives a `heartbeat` every 5 minutes, and it has details on whether there was a `click` in those 5 minutes, and if yes, the times the click happened.

## Demo

I created a tracker with `Mouse-click`, `Mouse-movement` and `Screen-Change` handlers activated. The `heartbeat` frequency was set to 12 seconds, i.e. every 12 seconds I received a `heartbeat` which mentioned all activities that were captured.

```
INFO[2019-03-29T11:35:29-07:00] starting activity tracker with 12 second frequency ...

INFO[2019-03-29T11:35:31-07:00] tracker worker working                        method=activity-tracker
INFO[2019-03-29T11:35:33-07:00] tracker worker working                        method=activity-tracker
INFO[2019-03-29T11:35:35-07:00] tracker worker working                        method=activity-tracker
INFO[2019-03-29T11:35:37-07:00] tracker worker working                        method=activity-tracker
INFO[2019-03-29T11:35:39-07:00] tracker worker working                        method=activity-tracker
INFO[2019-03-29T11:35:41-07:00] tracker worker working                        method=activity-tracker

INFO[2019-03-29T11:35:41-07:00] activity detected in the last 12 seconds.    
INFO[2019-03-29T11:35:41-07:00] Activity type:                               
INFO[2019-03-29T11:35:41-07:00] activityType : mouse-click times: 2          
INFO[2019-03-29T11:35:41-07:00] activityType : cursor-move times: 6          
INFO[2019-03-29T11:35:41-07:00] activityType : screen-change times: 2   
```

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

##### Note: Setting the `frequency` to a very less value might result in the tracker being invoked very frequently. The lowest value possible is every 10 seconds for now. (anything below that will revert it to default of 60s). The maximum value is 300s (5 minutes).

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
