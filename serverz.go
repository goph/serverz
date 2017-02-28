// Package serverz is a web app toolkit to easily manage server based environments
package serverz

import (
	"context"
	"net"
	"sync"

	"github.com/Sirupsen/logrus"
)

// ServerManager manages multiple Servers' lifecycle
type ServerManager struct {
	logger logrus.FieldLogger
}

// NewServerManager creates a new ServerManager
func NewServerManager(logger logrus.FieldLogger) *ServerManager {
	return &ServerManager{logger}
}

// StartServer creates a server starter function which can be called as a goroutine
func (sm *ServerManager) StartServer(server Server, listener net.Listener) func(ch chan<- error) {
	logger := sm.logger

	if s, ok := server.(*NamedServer); ok {
		logger = logger.WithField("name", s.Name)
	}

	return func(ch chan<- error) {
		logger.WithField("addr", listener.Addr().String()).Info("Starting server")
		ch <- server.Serve(listener)
	}
}

// ListenAndStartServer creates a server starter function which listens on a port and can be called as a goroutine
func (sm *ServerManager) ListenAndStartServer(server Server, addr string) func(ch chan<- error) {
	logger := sm.logger

	if s, ok := server.(*NamedServer); ok {
		logger = logger.WithField("name", s.Name)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal(err)
	}

	logger.WithField("addr", listener.Addr().String()).Info("Listening on address")

	return sm.StartServer(server, listener)
}

// StopServer creates a server stopper function which can be called as a goroutine
func (sm *ServerManager) StopServer(server Server, wg *sync.WaitGroup) func(ctx context.Context) {
	wg.Add(1)

	logger := sm.logger

	if s, ok := server.(*NamedServer); ok {
		logger = logger.WithField("name", s.Name)
	}

	return func(ctx context.Context) {
		logger.Info("Stopping server")

		err := server.Shutdown(ctx)
		if err != nil {
			sm.logger.Error(err)
		}

		wg.Done()
	}
}
