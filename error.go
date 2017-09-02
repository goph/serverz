package serverz

// multiError is returned by a composite closer.
type multiError []error

// Error implements the error interface.
func (e multiError) Error() string {
	return "Multiple errors happened"
}

// Errors implements the emperror.ErrorCollection interface.
func (e multiError) Errors() []error {
	return e
}

// ErrOrNil returns a nil when no errors are in the list.
func (e multiError) ErrOrNil() error {
	if len(e) == 0 {
		return nil
	}

	return e
}
