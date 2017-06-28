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
	Logger logger
}

// NewManager creates a new Manager.
func NewManager() *Manager {
	return &Manager{log.NewNopLogger()}
}

// StartServer creates a server starter function which can be called as a goroutine.
func (m *Manager) StartServer(server Server, lis net.Listener) func(ch chan<- error) {
	var name string
	if server, ok := server.(*NamedServer); ok {
		name = server.Name
	}

	return func(ch chan<- error) {
		level.Info(m.Logger).Log(
			"msg", "Starting server",
			"addr", lis.Addr().String(),
			"server", name,
		)
		ch <- server.Serve(lis)
	}
}

// ListenAndStartServer creates a server starter function which listens on a port and can be called as a goroutine
func (m *Manager) ListenAndStartServer(server Server, addr string) func(ch chan<- error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	var name string
	if server, ok := server.(*NamedServer); ok {
		name = server.Name
	}

	level.Info(m.Logger).Log(
		"msg", "Listening on address",
		"addr", addr,
		"server", name,
	)

	return m.StartServer(server, lis)
}

// StopServer creates a server stopper function which can be called as a goroutine
func (m *Manager) StopServer(server Server, wg *sync.WaitGroup) func(ctx context.Context) error {
	wg.Add(1)

	var name string
	if server, ok := server.(*NamedServer); ok {
		name = server.Name
	}

	return func(ctx context.Context) error {
		level.Info(m.Logger).Log(
			"msg", "Stopping server",
			"server", name,
		)

		err := server.Shutdown(ctx)

		wg.Done()

		return err
	}
}
