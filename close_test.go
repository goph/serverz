package serverz_test

import (
	"testing"

	"fmt"

	. "github.com/goph/serverz"
	"github.com/goph/serverz/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloserFunc_CallsUnderlyingFunc(t *testing.T) {
	var called bool

	closer := CloserFunc(func() {
		called = true
	})

	err := closer.Close()

	assert.NoError(t, err)

	assert.True(t, called)
}

func TestCloserFunc_RecoversErrorPanic(t *testing.T) {
	err := fmt.Errorf("internal error")

	closer := CloserFunc(func() {
		panic(err)
	})

	assert.EqualError(t, closer.Close(), "internal error")
}

func TestClosers(t *testing.T) {
	closer1 := new(mocks.Closer)
	closer1.On("Close").Return(nil)

	closer2 := new(mocks.Closer)
	closer2.On("Close").Return(nil)

	closer := Closers{closer1, closer2}

	assert.NoError(t, closer.Close())

	closer1.AssertExpectations(t)
	closer2.AssertExpectations(t)
}

func TestClosers_Empty(t *testing.T) {
	closer := Closers{}

	assert.NoError(t, closer.Close())
}

func TestClosers_Error(t *testing.T) {
	closer1 := new(mocks.Closer)
	closer1.On("Close").Return(nil)

	err := fmt.Errorf("error")
	closer2 := new(mocks.Closer)
	closer2.On("Close").Return(err)

	closer := Closers{closer1, closer2}

	merr := closer.Close()

	type errorCollection interface {
		Errors() []error
	}

	require.Error(t, merr)
	require.Implements(t, (*errorCollection)(nil), merr)
	assert.Contains(t, merr.(errorCollection).Errors(), err)

	closer1.AssertExpectations(t)
	closer2.AssertExpectations(t)
}

func TestClose(t *testing.T) {
	err := fmt.Errorf("error")
	closer := new(mocks.Closer)
	closer.On("Close").Return(err)

	assert.Equal(t, err, Close(closer))

	closer.AssertExpectations(t)
}
