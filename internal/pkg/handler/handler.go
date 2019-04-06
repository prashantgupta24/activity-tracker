package handler

import (
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
	"github.com/prashantgupta24/activity-tracker/pkg/system"
	log "github.com/sirupsen/logrus"
)

const (
	timeout = 100 //ms
)

//Instance is the main interface for a handler for the tracker
type Instance interface {
	Start(*log.Logger, chan *activity.Instance)
	Type() activity.Type
	Trigger(system.State) //used to activate pull-based handlers
	Close()
}
