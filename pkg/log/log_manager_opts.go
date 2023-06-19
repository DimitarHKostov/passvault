package log

const (
	defaultLogLevel = "info"
)

type LogOptsFn func(*LogOpts)

type LogOpts struct {
	Level string
}

func defaultLogOpts() LogOpts {
	return LogOpts{Level: defaultLogLevel}
}
