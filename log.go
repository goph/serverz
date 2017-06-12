package serverz

// logger needs to be satisfied by any loggers passed to the manager.
//
// See https://github.com/goph/log
type logger interface {
	WithField(key string, value interface{}) logger

	Info(args ...interface{})
	Error(args ...interface{})
}

// noopLogger is a default fallback implementation.
type noopLogger struct{}

func (l *noopLogger) WithField(key string, value interface{}) logger {
	return l
}

func (l *noopLogger) Info(args ...interface{}) {
}

func (l *noopLogger) Error(args ...interface{}) {
}
