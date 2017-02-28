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
	name := "Server"

	if s, ok := server.(*NamedServer); ok {
		name = s.Name
	}

	return func(ch chan<- error) {
		sm.logger.WithField("addr", listener.Addr().String()).Infof("%s started", name)
		ch <- server.Serve(listener)
	}
}

// ListenAndStartServer creates a server starter function which listens on a port and can be called as a goroutine
func (sm *ServerManager) ListenAndStartServer(server Server, addr string) func(ch chan<- error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		sm.logger.Fatal(err)
	}

	return sm.StartServer(server, listener)
}

// StopServer creates a server stopper function which can be called as a goroutine
func (sm *ServerManager) StopServer(server Server, wg *sync.WaitGroup) func(ctx context.Context) {
	wg.Add(1)

	name := "Server"

	if s, ok := server.(*NamedServer); ok {
		name = s.Name
	}

	return func(ctx context.Context) {
		sm.logger.Infof("Stopping %s", name)

		err := server.Shutdown(ctx)
		if err != nil {
			sm.logger.Error(err)
		}

		wg.Done()
	}
}
