package serverz

import "github.com/go-kit/kit/log"

// Option sets a value in an options instance.
type Option func(o *options)

// Logger returns a Option that sets the logger.
func Logger(l log.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}

// options holds a list of options frequently required by different components of the system.
type options struct {
	logger log.Logger
}

// newOptions returns a new options instance.
func newOptions(opt ...Option) *options {
	opts := new(options)

	for _, o := range opt {
		o(opts)
	}

	// Default logger
	if opts.logger == nil {
		opts.logger = log.NewNopLogger()
	}

	return opts
}

// Logger returns a log.Logger instance.
func (o *options) Logger() log.Logger {
	return o.logger
}