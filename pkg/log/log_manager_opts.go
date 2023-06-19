package log

import "passvault/pkg/types"

type LogOptsFn func(*LogOpts)

type LogOpts struct {
	Level string
}

func defaultLogOpts() LogOpts {
	return LogOpts{Level: types.DefaultLogLevel}
}
