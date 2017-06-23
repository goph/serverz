package serverz

import (
	"sync"
)

// multiError wraps multiple errors and exposes them as a single value.
type multiError struct {
	msg    string
	errors []error
	mu     *sync.RWMutex
}

// newMultiError returns a new multiError.
func newMultiError(msg string) *multiError {
	return &multiError{
		msg:    msg,
		errors: []error{},
		mu:     &sync.RWMutex{},
	}
}

// Error implements the error interface.
func (e *multiError) Error() string {
	return e.msg
}

// Errors implements the ErrorCollection interface from github.com/goph/stdlib/errors.
// This method is concurrent safe.
func (e *multiError) Errors() []error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.errors
}

// append allows to append an error to the list.
// This method is concurrent safe.
func (e *multiError) append(err error) {
	// Do not append nil errors.
	if err == nil {
		return
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.errors = append(e.errors, err)
}

// errOrNil spares an if in the caller code by returning nil when the error list is empty.
func (e *multiError) errOrNil() *multiError {
	if len(e.errors) == 0 {
		return nil
	}

	return e
}
