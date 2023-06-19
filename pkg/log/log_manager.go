package log

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type LogManager struct {
	logger *zap.Logger
}

func NewLogManager(logOpts ...LogOptsFn) *LogManager {
	opts := defaultLogOpts()

	for _, fn := range logOpts {
		fn(&opts)
	}

	var cfg zap.Config
	if err := json.Unmarshal(getConfig(opts.Level), &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

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

func getConfig(level string) []byte {
	configJSON := []byte(fmt.Sprintf(`{
		"level": "%s",
		"encoding": "json",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`, level))

	return configJSON
}
