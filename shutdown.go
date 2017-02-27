package serverz

import "github.com/Sirupsen/logrus"

// Shutdown manages an application shutdown by calling the registered handlers
type Shutdown struct {
	handlers []ShutdownHandler
	logger   logrus.FieldLogger
}

// ShutdownHandler is any function that has no parameters and can return an error
// Returned errors will be logged
type ShutdownHandler func() error

// ShutdownFunc wraps a function withot error return type
func ShutdownFunc(fn func()) ShutdownHandler {
	return func() error {
		fn()
		return nil
	}
}

// NewShutdown creates a new Shutdown
func NewShutdown(logger logrus.FieldLogger) *Shutdown {
	return &Shutdown{
		logger: logger,
	}
}

// Register appends new shutdown handlers to the list of existing ones
func (s *Shutdown) Register(handlers ...ShutdownHandler) {
	s.handlers = append(s.handlers, handlers...)
}

// RegisterAsFirst prepends new shutdown handlers to the list of existing ones
func (s *Shutdown) RegisterAsFirst(handlers ...ShutdownHandler) {
	s.handlers = append(handlers, s.handlers...)
}

// Handle is the panic recovery and shutdown handler
// It should be called as the last method in `main` (eg. using defer)
func (s *Shutdown) Handle() {
	// Try recovering from panic first
	v := recover()
	if v != nil {
		s.logger.Error(v)
	}

	s.logger.Info("Shutting down")

	// Loop through all the handlers and call them
	// Log any errors that may occur
	for _, handler := range s.handlers {
		err := handler()
		if err != nil {
			s.logger.Error(err)
		}
	}
}
