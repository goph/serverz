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
	lctx := logServerAddr(server, lis)

	return func(ch chan<- error) {
		m.Logger.Info("Starting server", lctx)
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

	lctx := logServerAddr(server, lis)

	m.Logger.Info("Listening on address", lctx)

	return m.StartServer(server, lis)
}

// StopServer creates a server stopper function which can be called as a goroutine
func (m *Manager) StopServer(server Server, wg *sync.WaitGroup) func(ctx context.Context) {
	wg.Add(1)

	lctx := logServer(server)

	return func(ctx context.Context) {
		m.Logger.Info("Stopping server", lctx)

		err := server.Shutdown(ctx)
		if err != nil && err != ctx.Err() {
			m.ErrorHandler.Handle(err)
		}

		wg.Done()
	}
}

func logServerAddr(s Server, l net.Listener) map[string]interface{} {
	ctx := logServer(s)

	ctx["addr"] = l.Addr().String()

	return ctx
}

func logServer(s Server) map[string]interface{} {
	ctx := make(map[string]interface{})

	if s, ok := s.(*NamedServer); ok {
		ctx["server"] = s.Name
	}

	return ctx
}
