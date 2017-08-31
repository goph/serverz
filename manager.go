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
			"addr", lis.Addr(),
			"server", getServerName(server),
		)

		ch <- server.Serve(lis)
	}
}

// ListenAndStartServer creates a server starter function which listens on a port and can be called as a goroutine.
func (m *Manager) ListenAndStartServer(server Server, addr net.Addr) (func(ch chan<- error), error) {
	var lis net.Listener
	var err error

	// Servers without address can still be managed via the Manager.
	if addr != nil {
		lis, err = net.Listen(addr.Network(), addr.String())
		if err != nil {
			return nil, err
		}
	} else {
		addr = NewAddr("none", "none")
		lis = newVirtualListener(addr)
	}

	level.Info(m.Logger).Log(
		"msg", "Listening on address",
		"addr", addr,
		"server", getServerName(server),
	)

	return m.StartServer(server, lis), nil
}

// StopServer creates a server stopper function which can be called as a goroutine.
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
