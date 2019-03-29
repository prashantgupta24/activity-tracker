# Activity tracker

[![codecov](https://codecov.io/gh/prashantgupta24/activity-tracker/branch/master/graph/badge.svg)](https://codecov.io/gh/prashantgupta24/activity-tracker) [![Go Report Card](https://goreportcard.com/badge/github.com/prashantgupta24/activity-tracker)](https://goreportcard.com/report/github.com/prashantgupta24/activity-tracker) [![version][version-badge]][RELEASES]

It is a libary that lets you monitor certain activities on your machine, and sends a heartbeat at a periodic (configurable) time detailing all the activity changes during that time. The activities that you want to monitor are **pluggable** services and can be added or removed according to your needs.

For example, you could have the tracker track Activities A, B and C, and you want it to monitor every 5 minutes. What it would do then is to send you heartbeats every 5 minutes, and each heartbeat would contain whether any of A, B or C occured within those 5 minutes, and if so, at what time.

As another example, let's say you want to monitor whether there was any mouse click on your machine and you want to be monitor every 5 minutes. What you do is start the `Activity Tracker` with just the `mouse click` service and `heartbeat` frequency set to 5 minutes. The `Start` function of the library gives you a channel which receives a `heartbeat` every 5 minutes, and it has details on whether there was a `click` in those 5 minutes.

## Heartbeat struct

It is the data packet sent from the tracker library to the user.

	type Heartbeat struct {
		WasAnyActivity bool //whether any activity was detected 		
		ActivityMap       map[*activity.Type][]time.Time //activity type with its times
		Time           time.Time                    //heartbeat time
	}

`WasAnyActivity` tells if there was any activity within that time frame
If there was, then the ActivityMap will tell you what type of activity
it was and what times it occured.

The `Time` field is the time of the Heartbeat sent (not to be confused with
the activity time, which is the time the activity occured within the time frame)

## Pluggable services/activities

In the example above, `mouse click` was just one of the activities currently supported. Any new type of activity can be easily added, it just needs to implement the `Service` interface and define what type of activity it is going to track, that's it! It can be plugged in with a tracker and then the tracker will include those activity checks in its heartbeat.

## Currently supported list of services

- Mouse click (whether any mouse click happened during the time frame)
- Mouse cursor movement (whether the mouse cursor was moved during the time frame)
- Screen change (whether the screen was changed anytime within that time frame)

## Example output

I created a tracker with `

[version-badge]: https://img.shields.io/github/release/prashantgupta24/activity-tracker.svg
[RELEASES]: https://github.com/prashantgupta24/activity-tracker/releases
