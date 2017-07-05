package ext

import "github.com/goph/serverz/internal/errors"

// NoopCloser is a dummy Closer implementation which can be used as a fallback.
var NoopCloser = CloserFunc(func() {})

// Closer is an alias interface to io.Closer.
type Closer interface {
	Close() error
}

// CloserFunc makes any function a Closer.
type CloserFunc func()

// Close calls the underlying function and converts any panic to an error.
func (f CloserFunc) Close() (err error) {
	defer func() {
		err = errors.Recover(recover())
	}()

	f()

	return err
}

// Closers is a collection of Closer instances.
type Closers []Closer

// Close calls the underlying Closer instances and returns all their errors as a single value.
func (c Closers) Close() error {
	if len(c) == 0 {
		return nil
	}

	errBuilder := errors.NewMultiErrorBuilder()

	for _, closer := range c {
		err := closer.Close()

		errBuilder.Add(err)
	}

	return errBuilder.ErrOrNil()
}
