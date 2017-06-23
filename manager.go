package serverz

import (
	"context"
	"net"
	"sync"
)

// Manager manages multiple Servers' lifecycle.
type Manager struct {
	Logger       logger
	ErrorHandler errorHandler
}

// NewManager creates a new Manager.
func NewManager() *Manager {
	return &Manager{
		&noopLogger{},
		&noopErrorHandler{},
	}
}

// StartServer creates a server starter function which can be called as a goroutine.
func (m *Manager) StartServer(server Server, lis net.Listener) func(ch chan<- error) {
	var name string
	if server, ok := server.(*NamedServer); ok {
		name = server.Name
	}

	return func(ch chan<- error) {
		m.Logger.Log("level", "info", "msg", "Starting server", "addr", lis.Addr().String(), "server", name)
		ch <- server.Serve(lis)
	}
}

// ListenAndStartServer creates a server starter function which listens on a port and can be called as a goroutine
func (m *Manager) ListenAndStartServer(server Server, addr string) func(ch chan<- error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		m.ErrorHandler.Handle(err)
		panic(err)
	}

	var name string
	if server, ok := server.(*NamedServer); ok {
		name = server.Name
	}

	m.Logger.Log("level", "info", "msg", "Listening on address", "addr", addr, "server", name)

	return m.StartServer(server, lis)
}

// StopServer creates a server stopper function which can be called as a goroutine
func (m *Manager) StopServer(server Server, wg *sync.WaitGroup) func(ctx context.Context) {
	wg.Add(1)

	var name string
	if server, ok := server.(*NamedServer); ok {
		name = server.Name
	}

	return func(ctx context.Context) {
		m.Logger.Log("level", "info", "msg", "Stopping server", "server", name)

		err := server.Shutdown(ctx)
		if err != nil && err != ctx.Err() {
			m.ErrorHandler.Handle(err)
		}

		wg.Done()
	}
}
