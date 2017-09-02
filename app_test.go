package serverz_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	. "github.com/goph/serverz"
	"github.com/goph/serverz/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAppServerIsAServer(t *testing.T) {
	assert.Implements(t, (*Server)(nil), new(AppServer))
	assert.Implements(t, (*ListenServer)(nil), new(AppServer))
}

func TestAppServer_Serve(t *testing.T) {
	serverMock := new(mocks.Server)
	serverMock.On("Serve", mock.Anything).Return(nil)

	buf := new(bytes.Buffer)
	logger := log.NewLogfmtLogger(buf)

	server := &AppServer{
		Server: serverMock,

		Name: "server",

		Logger: logger,
	}

	lis := NewListener(NewAddr("none", "none"))

	server.Serve(lis)

	assert.Equal(t, "msg=\"Starting server\" addr=none server=server\n", buf.String())

	serverMock.AssertExpectations(t)
}

func TestAppServer_ListenAndServe(t *testing.T) {
	serverMock := new(mocks.Server)
	serverMock.On("Serve", mock.Anything).Return(nil)

	buf := new(bytes.Buffer)
	logger := log.NewLogfmtLogger(buf)

	server := &AppServer{
		Server: serverMock,

		Name: "server",

		Logger: logger,
	}

	server.ListenAndServe(nil)

	assert.Equal(t, "msg=\"Listening on address\" addr=none server=server\nmsg=\"Starting server\" addr=none server=server\n", buf.String())

	serverMock.AssertExpectations(t)
}

func TestAppServer_ListenAndServe_Real(t *testing.T) {
	serverMock := new(mocks.Server)
	serverMock.On("Serve", mock.Anything).Return(nil)

	server := &AppServer{
		Server: serverMock,

		Logger: log.NewNopLogger(),
	}

	server.ListenAndServe(NewAddr("tcp", ":0"))

	serverMock.AssertExpectations(t)
}

func TestAppServer_Shutdown(t *testing.T) {
	ctx := context.Background()

	serverMock := new(mocks.Server)
	serverMock.On("Shutdown", ctx).Return(nil)

	buf := new(bytes.Buffer)
	logger := log.NewLogfmtLogger(buf)

	server := &AppServer{
		Server: serverMock,

		Name: "server",

		Logger: logger,
	}

	server.Shutdown(ctx)

	assert.Equal(t, "msg=\"Attempting to shut server gracefully down\" server=server\n", buf.String())

	serverMock.AssertExpectations(t)
}

func TestAppServer_Close(t *testing.T) {
	serverMock := new(mocks.Server)
	serverMock.On("Close").Return(nil)

	closerMock := new(mocks.Closer)
	closerMock.On("Close").Return(nil)

	buf := new(bytes.Buffer)
	logger := log.NewLogfmtLogger(buf)

	server := &AppServer{
		Server: serverMock,

		Name:   "server",
		Closer: closerMock,

		Logger: logger,
	}

	server.Close()

	assert.Equal(t, "msg=\"Closing server\" server=server\n", buf.String())

	serverMock.AssertExpectations(t)
	closerMock.AssertExpectations(t)
}
