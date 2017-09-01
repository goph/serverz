package serverz_test

import (
	"testing"

	"context"
	"net"
	"sync"

	. "github.com/goph/serverz"
	"github.com/goph/serverz/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestManager_StartServer(t *testing.T) {
	serverManager := NewManager()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer lis.Close()

	server := &mocks.Server{}
	server.On("Serve", lis).Return(nil)

	ch := make(chan error, 1)
	serverManager.StartServer(server, lis)(ch)

	server.AssertExpectations(t)
}

func TestManager_ListenAndStartServer(t *testing.T) {
	serverManager := NewManager()

	server := &mocks.Server{}
	server.On("Serve", mock.Anything).Return(func(lis net.Listener) error {
		lis.Close()

		return nil
	})

	addr := NewAddr("tcp", "127.0.0.1:0")
	ch := make(chan error, 1)
	starter, err := serverManager.ListenAndStartServer(server, addr)

	starter(ch)

	require.NoError(t, err)

	server.AssertExpectations(t)
}

func TestManager_ListenAndStartServer_NoAddr(t *testing.T) {
	serverManager := NewManager()

	server := &mocks.Server{}
	server.On("Serve", mock.Anything).Return(nil)

	ch := make(chan error, 1)
	starter, err := serverManager.ListenAndStartServer(server, nil)

	starter(ch)

	require.NoError(t, err)

	server.AssertExpectations(t)
}

func TestManager_StopServer(t *testing.T) {
	serverManager := NewManager()

	ctx := context.Background()
	server := &mocks.Server{}
	server.On("Shutdown", ctx).Return(nil)

	var wg sync.WaitGroup

	stopper := serverManager.StopServer(server, &wg)

	stopper(ctx)

	server.AssertExpectations(t)
}
