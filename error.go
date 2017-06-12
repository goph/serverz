package serverz

// errorHandler handles an error.
type errorHandler interface {
	Handle(err error)
}

// noopErrorHandler is a default fallback in case there is no error handler configured.
type noopErrorHandler struct{}

func (e *noopErrorHandler) Handle(err error) {}
