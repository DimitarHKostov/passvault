package log

import "go.uber.org/zap"

var (
	logManager *LogManager
)

type LogManager struct {
	Logger *zap.Logger
}

func Get() *LogManager {
	if logManager == nil {
		logger, _ := zap.NewProduction()
		logManager = &LogManager{Logger: logger}
	}

	return logManager
}
