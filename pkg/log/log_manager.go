package log

import "go.uber.org/zap"

type LogManager struct {
	logger *zap.Logger
}

func NewLogManager() *LogManager {
	logger, _ := zap.NewProduction()
	logManager := &LogManager{logger: logger}

	return logManager
}

func (lm *LogManager) LogInfo(message string) {
	lm.logger.Info(message)
}

func (lm *LogManager) LogError(message string) {
	lm.logger.Error(message)
}

func (lm *LogManager) LogDebug(message string) {
	lm.logger.Debug(message)
}

func (lm *LogManager) LogFatal(message string) {
	lm.logger.Fatal(message)
}

func (lm *LogManager) LogPanic(message string) {
	lm.logger.Panic(message)
}
