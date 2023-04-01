package log

type LogManagerInterface interface {
	LogInfo(string)
	LogError(string)
	LogDebug(string)
	LogFatal(string)
	LogPanic(string)
}
