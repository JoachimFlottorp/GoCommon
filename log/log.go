package log

import (
	"go.uber.org/zap"
)

var (
	DefaultConfig = zap.Config{
		Level: zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding: "console",
	}
)

func InitLogger(conf zap.Config) error {
	logger, err := conf.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)

	return nil
}