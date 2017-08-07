package aio_test

import (
	"net"
	"testing"

	"github.com/goph/serverz"
	"github.com/goph/serverz/aio"
	"github.com/goph/serverz/internal"
	"github.com/goph/serverz/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestServerIsAServer(t *testing.T) {
	assert.Implements(t, (*serverz.Server)(nil), new(aio.Server))
}

func TestServer_Name(t *testing.T) {
	server := &aio.Server{
		Name: "name",
	}

	assert.Implements(t, (*internal.NamedServer)(nil), server)
	assert.Equal(t, "name", server.GetName())
}

func TestServer_Addr(t *testing.T) {
	addr := &net.TCPAddr{}
	server := &aio.Server{
		Addr: addr,
	}

	assert.Implements(t, (*serverz.AddrServer)(nil), server)
	assert.Equal(t, addr, server.GetAddr())
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
