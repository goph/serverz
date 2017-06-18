package serverz

// logger needs to be satisfied by any loggers passed to the manager.
//
// See https://github.com/goph/log
type logger interface {
	Info(args ...interface{})
}

// noopLogger is a default fallback implementation.
type noopLogger struct{}

func (l *noopLogger) Info(args ...interface{}) {
}
