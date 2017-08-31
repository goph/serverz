package serverz

import (
	"context"
	"net"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Manager manages multiple Servers' lifecycle.
type Manager struct {
	Logger log.Logger
}

// NewManager creates a new Manager.
func NewManager() *Manager {
	return &Manager{log.NewNopLogger()}
}

// StartServer creates a server starter function which can be called as a goroutine.
func (m *Manager) StartServer(server Server, lis net.Listener) func(ch chan<- error) {
	return func(ch chan<- error) {
		level.Info(m.Logger).Log(
			"msg", "Starting server",
			"addr", lis.Addr().String(),
			"server", getServerName(server),
		)
		ch <- server.Serve(lis)
	}
}

// ListenAndStartServer creates a server starter function which listens on a port and can be called as a goroutine
func (m *Manager) ListenAndStartServer(server Server, network string, addr string) func(ch chan<- error) {
	lis, err := net.Listen(network, addr)
	if err != nil {
		panic(err)
	}

	level.Info(m.Logger).Log(
		"msg", "Listening on address",
		"addr", addr,
		"server", getServerName(server),
	)

	return m.StartServer(server, lis)
}

// StopServer creates a server stopper function which can be called as a goroutine
func (m *Manager) StopServer(server Server, wg *sync.WaitGroup) func(ctx context.Context) error {
	wg.Add(1)

	return func(ctx context.Context) error {
		level.Info(m.Logger).Log(
			"msg", "Stopping server",
			"server", getServerName(server),
		)

		err := server.Shutdown(ctx)

		wg.Done()

		return err
	}
}
