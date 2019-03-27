# Activity tracker

[![codecov](https://codecov.io/gh/prashantgupta24/activity-tracker/branch/master/graph/badge.svg)](https://codecov.io/gh/prashantgupta24/activity-tracker) [![Go Report Card](https://goreportcard.com/badge/github.com/prashantgupta24/activity-tracker)](https://goreportcard.com/report/github.com/prashantgupta24/activity-tracker)

It is a libary that lets you monitor certain activities on your machine, and sends a heartbeat at a periodic (configurable) time detailing all the activity changes during that time. The activities that you want to monitor are pluggable services and can be added or removed according to your needs.

For example, let's say you want to monitor whether there was any mouse click on your machine and you want to be notified every 5 minutes if there was one or not. What you do is start the `Activity Tracker` with just the `mouse click` service and `heartbeat` frequency set to 5 minutes. The `Start` function of the library gives you a channel which receives a `heartbeat` every 5 minutes, and it has details on whether there was a `click` in those 5 minutes.
