package aio_test

import (
	"testing"

	"github.com/goph/serverz"
	"github.com/goph/serverz/aio"
	"github.com/goph/serverz/mocks"
	"github.com/stretchr/testify/assert"
)

func TestServerIsAServer(t *testing.T) {
	assert.Implements(t, (*serverz.Server)(nil), new(aio.Server))
}

func TestServer_Name(t *testing.T) {
	server := &aio.Server{
		ServerName: "name",
	}

	assert.Equal(t, "name", server.Name())
}

func TestServer_Close(t *testing.T) {
	s := &mocks.Server{}
	s.On("Close").Return(nil)

	c := &mocks.Closer{}
	c.On("Close").Return(nil)

	server := &aio.Server{
		Server: s,
		Closer: c,
	}

	server.Close()

	s.AssertCalled(t, "Close")
	c.AssertCalled(t, "Close")
	s.AssertExpectations(t)
}
