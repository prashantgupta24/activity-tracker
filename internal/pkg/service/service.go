package service

import (
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
	log "github.com/sirupsen/logrus"
)

const (
	timeout = 100 //ms
)

type Instance interface {
	Start(*log.Logger, chan *activity.Type)
	Type() activity.Type
	Trigger()
	Close()
}
