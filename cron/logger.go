// Wraps Cron logging in zap's SugaredLogger if desired.

package cron

import (
	"github.com/JoachimFlottorp/GoCommon/log"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type l struct{}

func WithLogger(shouldLog bool) cron.Logger {
	log.InitLogger(log.DefaultConfig)
	
	if shouldLog {
		return l{}
	}
	return cron.DefaultLogger
}

func (l) Info(msg string, keysAndValues ...interface{}) {
	zap.S().Infof(msg, keysAndValues...)
}

func (l) Error(err error, msg string, keysAndValues ...interface{}) {
	zap.S().Errorf(msg, err, keysAndValues)
}