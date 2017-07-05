package serverz_test

import (
	"testing"

	"context"
	"net"
	"sync"

	"github.com/goph/serverz"
	"github.com/goph/serverz/mocks"
	"github.com/stretchr/testify/mock"
)

type testListener struct{}

func (l *testListener) Accept() (net.Conn, error) {
	return nil, nil
}

func (l *testListener) Close() error {
	return nil
}

func (l *testListener) Addr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 65123,
	}
}

func TestServerManagerStartServer(t *testing.T) {
	serverManager := serverz.NewManager()

	lis := &testListener{}
	server := &mocks.Server{}
	server.On("Serve", lis).Return(nil)

	ch := make(chan error, 1)
	serverManager.StartServer(server, lis)(ch)
	close(ch)

	server.AssertCalled(t, "Serve", lis)
	server.AssertExpectations(t)
}

func TestServerManagerListenAndStartServer(t *testing.T) {
	serverManager := serverz.NewManager()

	server := &mocks.Server{}
	server.On("Serve", mock.Anything).Return(func(list net.Listener) error {
		list.Close()

		return nil
	})

	ch := make(chan error, 1)
	serverManager.ListenAndStartServer(server, "127.0.0.1:65123")(ch)
	close(ch)

	server.AssertCalled(t, "Serve", mock.Anything)
	server.AssertExpectations(t)
}

func TestServerManagerStopServer(t *testing.T) {
	serverManager := serverz.NewManager()

	ctx := context.Background()
	server := &mocks.Server{}
	server.On("Shutdown", ctx).Return(nil)

	var wg sync.WaitGroup

	f := serverManager.StopServer(server, &wg)

	f(ctx)

	server.AssertCalled(t, "Shutdown", ctx)
	server.AssertExpectations(t)
}
