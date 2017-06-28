package serverz

// logger needs to be satisfied by any loggers passed to the manager.
//
// See https://github.com/goph/log
type logger interface {
	Log(keyvals ...interface{}) error
}
