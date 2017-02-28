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

	testStartServer(t, server)
}

func TestServerManagerStartServer_NamedServer(t *testing.T) {
	server := &serverz.NamedServer{
		Server: &TestServer{},
		Name:   "ServerName",
	}

	testStartServer(t, server)
}

func testStartServer(t *testing.T, server serverz.Server) {
	logger, hook := logrus_test.NewNullLogger()
	serverManager := serverz.NewServerManager(logger)

	ch := make(chan error, 1)
	serverManager.StartServer(server, &TestListener{})(ch)
	close(ch)

	var called bool

	if s, ok := server.(*serverz.NamedServer); ok {
		called = s.Server.(*TestServer).ServeCalled
	} else {
		called = server.(*TestServer).ServeCalled
	}

	if got, want := called, true; got != want {
		t.Fatal("serverz: Server.Serve should be called")
	}

	if got, want := hook.LastEntry().Message, "Starting server"; got != want {
		t.Fatal("serverz: server should be started")
	}

	if got, want := hook.LastEntry().Level, logrus.InfoLevel; got != want {
		t.Fatal("serverz: server start is info level")
	}

	if got, want := hook.LastEntry().Data["addr"], "127.0.0.1:65123"; got != want {
		t.Fatal("serverz: address should be logged")
	}

	if s, ok := server.(*serverz.NamedServer); ok {
		if got, want := hook.LastEntry().Data["server"], s.Name; got != want {
			t.Fatal("serverz: server should log it's name")
		}
	}
}

func TestServerManagerListenAndStartServer(t *testing.T) {
	server := &TestServer{}

	testListenAndStartServer(t, server)
}

func TestServerManagerListenAndStartServer_NamedServer(t *testing.T) {
	server := &serverz.NamedServer{
		Server: &TestServer{},
		Name:   "ServerName",
	}

	testListenAndStartServer(t, server)
}

func testListenAndStartServer(t *testing.T, server serverz.Server) {
	logger, hook := logrus_test.NewNullLogger()
	serverManager := serverz.NewServerManager(logger)

	ch := make(chan error, 1)
	serverManager.ListenAndStartServer(server, "127.0.0.1:65123")(ch)
	server.Close()
	close(ch)

	var called bool

	if s, ok := server.(*serverz.NamedServer); ok {
		called = s.Server.(*TestServer).ServeCalled
	} else {
		called = server.(*TestServer).ServeCalled
	}

	if got, want := called, true; got != want {
		t.Fatal("serverz: Server.Serve should be called")
	}

	if got, want := hook.Entries[0].Message, "Listening on address"; got != want {
		t.Fatal("serverz: server should listen on address")
	}

	if got, want := hook.Entries[0].Level, logrus.InfoLevel; got != want {
		t.Fatal("serverz: listen is info level")
	}

	if got, want := hook.Entries[0].Data["addr"], "127.0.0.1:65123"; got != want {
		t.Fatal("serverz: server should log it's addr")
	}

	if got, want := hook.LastEntry().Message, "Starting server"; got != want {
		t.Fatal("serverz: server should be started")
	}

	if got, want := hook.LastEntry().Level, logrus.InfoLevel; got != want {
		t.Fatal("serverz: server start is info level")
	}

	if got, want := hook.LastEntry().Data["addr"], "127.0.0.1:65123"; got != want {
		t.Fatal("serverz: address should be logged")
	}

	if s, ok := server.(*serverz.NamedServer); ok {
		if got, want := hook.LastEntry().Data["server"], s.Name; got != want {
			t.Fatal("serverz: server should log it's name")
		}
	}
}

func TestServerManagerStopServer(t *testing.T) {
	server := &TestServer{}

	testStopServer(t, server)
}

func TestServerManagerStopServer_NamedServer(t *testing.T) {
	server := &serverz.NamedServer{
		Server: &TestServer{},
		Name:   "ServerName",
	}

	testStopServer(t, server)
}

func testStopServer(t *testing.T, server serverz.Server) {
	logger, hook := logrus_test.NewNullLogger()
	serverManager := serverz.NewServerManager(logger)

	var wg sync.WaitGroup

	f := serverManager.StopServer(server, &wg)

	ctx := context.Background()
	f(ctx)

	var called bool

	if s, ok := server.(*serverz.NamedServer); ok {
		called = s.Server.(*TestServer).ShutdownCalled
	} else {
		called = server.(*TestServer).ShutdownCalled
	}

	if got, want := called, true; got != want {
		t.Fatal("serverz: Server.Shutdown should be called")
	}

	if got, want := hook.LastEntry().Message, "Stopping server"; got != want {
		t.Fatal("serverz: server should be stopped")
	}

	if got, want := hook.LastEntry().Level, logrus.InfoLevel; got != want {
		t.Fatal("serverz: server stop is info level")
	}

	if s, ok := server.(*serverz.NamedServer); ok {
		if got, want := hook.LastEntry().Data["server"], s.Name; got != want {
			t.Fatal("serverz: server should log it's name")
		}
	}
}
