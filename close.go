package serverz

import (
	"errors"
	"fmt"
)

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
		r := recover()

		if r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = fmt.Errorf("Unknown panic, received: %v", r)
			}
		}
	}()

	f()

	return
}

// Closers is a collection of Closer instances.
type Closers []Closer

// Close calls the underlying Closer instances and returns all their errors as a single value.
func (c Closers) Close() error {
	if len(c) == 0 {
		return nil
	}

	var merr multiError

	for _, closer := range c {
		err := closer.Close()

		if err != nil {
			merr = append(merr, err)
		}
	}

	return merr.ErrOrNil()
}

// Close calls Close method on a struct if it implements the Closer interface.
func Close(closer interface{}) error {
	if c, ok := closer.(Closer); ok {
		return c.Close()
	}

	return nil
}
