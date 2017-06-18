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
	l := m.Logger

	if s, ok := server.(*NamedServer); ok {
		l = l.WithField("server", s.Name).(logger)
	}

	return func(ch chan<- error) {
		l.WithField("addr", lis.Addr().String()).(logger).Info("Starting server")
		ch <- server.Serve(lis)
	}
}

// ListenAndStartServer creates a server starter function which listens on a port and can be called as a goroutine
func (m *Manager) ListenAndStartServer(server Server, addr string) func(ch chan<- error) {
	l := m.Logger

	if s, ok := server.(*NamedServer); ok {
		l = l.WithField("server", s.Name).(logger)
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		m.ErrorHandler.Handle(err)
		panic(err)
	}

	l.WithField("addr", lis.Addr().String()).(logger).Info("Listening on address")

	return m.StartServer(server, lis)
}

// StopServer creates a server stopper function which can be called as a goroutine
func (m *Manager) StopServer(server Server, wg *sync.WaitGroup) func(ctx context.Context) {
	wg.Add(1)

	l := m.Logger

	if s, ok := server.(*NamedServer); ok {
		l = l.WithField("server", s.Name).(logger)
	}

	return func(ctx context.Context) {
		l.Info("Stopping server")

		err := server.Shutdown(ctx)
		if err != nil && err != ctx.Err() {
			m.ErrorHandler.Handle(err)
		}

		wg.Done()
	}
}
