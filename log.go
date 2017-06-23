package serverz

// logger needs to be satisfied by any loggers passed to the manager.
//
// See https://github.com/goph/log
type logger interface {
	Log(keyvals ...interface{}) error
}

// noopLogger is a default fallback implementation.
type noopLogger struct{}

func (l *noopLogger) Log(keyvals ...interface{}) error {
	return nil
}
