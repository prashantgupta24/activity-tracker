package service

import "github.com/prashantgupta24/activity-tracker/pkg/activity"

const (
	timeout = 100 //ms
)

type Instance interface {
	Start(chan *activity.Type) chan struct{}
}
