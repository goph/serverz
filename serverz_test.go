package serverz_test

import (
	"context"
	"net"
	"testing"

	"sync"

	"github.com/Sirupsen/logrus"
	logrus_test "github.com/Sirupsen/logrus/hooks/test"
	"github.com/sagikazarmark/serverz"
)

type TestListener struct{}

func (l *TestListener) Accept() (net.Conn, error) {
	return nil, nil
}

func (l *TestListener) Close() error {
	return nil
}

func (l *TestListener) Addr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 65123,
	}
}

type TestServer struct {
	listener net.Listener

	ServeCalled    bool
	ShutdownCalled bool
	CloseCalled    bool

	ServeError    error
	ShutdownError error
	CloseError    error
}

func (s *TestServer) Serve(l net.Listener) error {
	s.listener = l
	s.ServeCalled = true

	return s.ServeError
}

func (s *TestServer) Shutdown(ctx context.Context) error {
	s.ShutdownCalled = true

	return s.ShutdownError
}

func (s *TestServer) Close() error {
	if s.listener != nil {
		s.listener.Close()
	}

	s.CloseCalled = true

	return s.CloseError
}

func TestServerManagerStartServer(t *testing.T) {
	server := &TestServer{}
	logger, hook := logrus_test.NewNullLogger()
	serverManager := serverz.NewServerManager(logger)

	f := serverManager.StartServer(server, &TestListener{})

	ch := make(chan error, 1)
	f(ch)

	if got, want := server.ServeCalled, true; got != want {
		t.Fatal("serverz: Server.Serve should be called")
	}

	if got, want, gotLevel, wantLevel := hook.LastEntry().Message, "Server started", hook.LastEntry().Level, logrus.InfoLevel; got != want || gotLevel != wantLevel {
		t.Fatal("serverz: server should be started")
	}

	if got, want := hook.LastEntry().Data["addr"], "127.0.0.1:65123"; got != want {
		t.Fatal("serverz: address should be logged")
	}
}

func TestServerManagerStartServer_NamedServer(t *testing.T) {
	server := &serverz.NamedServer{
		Server: &TestServer{},
		Name:   "ServerName",
	}

	logger, hook := logrus_test.NewNullLogger()
	serverManager := serverz.NewServerManager(logger)

	f := serverManager.StartServer(server, &TestListener{})

	ch := make(chan error, 1)
	f(ch)

	if got, want := server.Server.(*TestServer).ServeCalled, true; got != want {
		t.Fatal("serverz: Server.Serve should be called")
	}

	if got, want, gotLevel, wantLevel := hook.LastEntry().Message, "ServerName started", hook.LastEntry().Level, logrus.InfoLevel; got != want || gotLevel != wantLevel {
		t.Fatal("serverz: server should be started")
	}

	if got, want := hook.LastEntry().Data["addr"], "127.0.0.1:65123"; got != want {
		t.Fatal("serverz: address should be logged")
	}
}

func TestServerManagerListenAndStartServer(t *testing.T) {
	server := &TestServer{}

	logger, hook := logrus_test.NewNullLogger()
	serverManager := serverz.NewServerManager(logger)

	f := serverManager.ListenAndStartServer(server, "127.0.0.1:65123")

	ch := make(chan error, 1)
	f(ch)
	server.Close()

	if got, want := server.ServeCalled, true; got != want {
		t.Fatal("serverz: Server.Serve should be called")
	}

	if got, want, gotLevel, wantLevel := hook.LastEntry().Message, "Server started", hook.LastEntry().Level, logrus.InfoLevel; got != want || gotLevel != wantLevel {
		t.Fatal("serverz: server should be started")
	}

	if got, want := hook.LastEntry().Data["addr"], "127.0.0.1:65123"; got != want {
		t.Fatal("serverz: address should be logged")
	}
}

func TestServerManagerListenAndStartServer_NamedServer(t *testing.T) {
	server := &serverz.NamedServer{
		Server: &TestServer{},
		Name:   "ServerName",
	}

	logger, hook := logrus_test.NewNullLogger()
	serverManager := serverz.NewServerManager(logger)

	f := serverManager.ListenAndStartServer(server, "127.0.0.1:65123")

	ch := make(chan error, 1)
	f(ch)
	server.Close()

	if got, want := server.Server.(*TestServer).ServeCalled, true; got != want {
		t.Fatal("serverz: Server.Serve should be called")
	}

	if got, want, gotLevel, wantLevel := hook.LastEntry().Message, "ServerName started", hook.LastEntry().Level, logrus.InfoLevel; got != want || gotLevel != wantLevel {
		t.Fatal("serverz: server should be started")
	}

	if got, want := hook.LastEntry().Data["addr"], "127.0.0.1:65123"; got != want {
		t.Fatal("serverz: address should be logged")
	}
}

func TestServerManagerStopServer(t *testing.T) {
	server := &TestServer{}

	logger, hook := logrus_test.NewNullLogger()
	serverManager := serverz.NewServerManager(logger)

	var wg sync.WaitGroup

	f := serverManager.StopServer(server, &wg)

	ctx := context.Background()
	f(ctx)

	if got, want := server.ShutdownCalled, true; got != want {
		t.Fatal("serverz: Server.Shutdown should be called")
	}

	if got, want, gotLevel, wantLevel := hook.LastEntry().Message, "Stopping Server", hook.LastEntry().Level, logrus.InfoLevel; got != want || gotLevel != wantLevel {
		t.Fatal("serverz: server should be started")
	}
}

func TestServerManagerStopServer_NamedServer(t *testing.T) {
	server := &serverz.NamedServer{
		Server: &TestServer{},
		Name:   "ServerName",
	}

	logger, hook := logrus_test.NewNullLogger()
	serverManager := serverz.NewServerManager(logger)

	var wg sync.WaitGroup

	f := serverManager.StopServer(server, &wg)

	ctx := context.Background()
	f(ctx)

	if got, want := server.Server.(*TestServer).ShutdownCalled, true; got != want {
		t.Fatal("serverz: Server.Shutdown should be called")
	}

	if got, want, gotLevel, wantLevel := hook.LastEntry().Message, "Stopping ServerName", hook.LastEntry().Level, logrus.InfoLevel; got != want || gotLevel != wantLevel {
		t.Fatal("serverz: server should be started")
	}
}
