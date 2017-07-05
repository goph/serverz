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

func TestServerManagerStartServer(t *testing.T) {
	serverManager := serverz.NewManager()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer lis.Close()

	server := &mocks.Server{}
	server.On("Serve", lis).Return(nil)

	ch := make(chan error, 1)
	serverManager.StartServer(server, lis)(ch)

	server.AssertCalled(t, "Serve", lis)
	server.AssertExpectations(t)
}

func TestServerManagerListenAndStartServer(t *testing.T) {
	serverManager := serverz.NewManager()

	server := &mocks.Server{}
	server.On("Serve", mock.Anything).Return(func(lis net.Listener) error {
		lis.Close()

		return nil
	})

	ch := make(chan error, 1)
	serverManager.ListenAndStartServer(server, "127.0.0.1:0")(ch)

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
