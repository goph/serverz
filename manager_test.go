package serverz_test

import (
	"testing"

	"context"
	"net"
	"sync"

	"github.com/goph/serverz"
)

type testListener struct{}

func (l *testListener) Accept() (net.Conn, error) {
	return nil, nil
}

func (l *testListener) Close() error {
	return nil
}

func (l *testListener) Addr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 65123,
	}
}

type testServer struct {
	listener net.Listener

	ServeCalled    bool
	ShutdownCalled bool
	CloseCalled    bool

	ServeError    error
	ShutdownError error
	CloseError    error
}

func (s *testServer) Serve(l net.Listener) error {
	s.listener = l
	s.ServeCalled = true

	return s.ServeError
}

func (s *testServer) Shutdown(ctx context.Context) error {
	s.ShutdownCalled = true

	return s.ShutdownError
}

func (s *testServer) Close() error {
	if s.listener != nil {
		s.listener.Close()
	}

	s.CloseCalled = true

	return s.CloseError
}

func TestServerManagerStartServer(t *testing.T) {
	server := &testServer{}

	testStartServer(t, server)
}

func TestServerManagerStartServer_NamedServer(t *testing.T) {
	server := &serverz.NamedServer{
		Server: &testServer{},
		Name:   "ServerName",
	}

	testStartServer(t, server)
}

func testStartServer(t *testing.T, server serverz.Server) {
	serverManager := serverz.NewManager()

	ch := make(chan error, 1)
	serverManager.StartServer(server, &testListener{})(ch)
	close(ch)

	var called bool

	if s, ok := server.(*serverz.NamedServer); ok {
		called = s.Server.(*testServer).ServeCalled
	} else {
		called = server.(*testServer).ServeCalled
	}

	if got, want := called, true; got != want {
		t.Fatal("serverz: Server.Serve should be called")
	}
}

func TestServerManagerListenAndStartServer(t *testing.T) {
	server := &testServer{}

	testListenAndStartServer(t, server)
}

func TestServerManagerListenAndStartServer_NamedServer(t *testing.T) {
	server := &serverz.NamedServer{
		Server: &testServer{},
		Name:   "ServerName",
	}

	testListenAndStartServer(t, server)
}

func testListenAndStartServer(t *testing.T, server serverz.Server) {
	serverManager := serverz.NewManager()

	ch := make(chan error, 1)
	serverManager.ListenAndStartServer(server, "127.0.0.1:65123")(ch)
	server.Close()
	close(ch)

	var called bool

	if s, ok := server.(*serverz.NamedServer); ok {
		called = s.Server.(*testServer).ServeCalled
	} else {
		called = server.(*testServer).ServeCalled
	}

	if got, want := called, true; got != want {
		t.Fatal("serverz: Server.Serve should be called")
	}
}

func TestServerManagerStopServer(t *testing.T) {
	server := &testServer{}

	testStopServer(t, server)
}

func TestServerManagerStopServer_NamedServer(t *testing.T) {
	server := &serverz.NamedServer{
		Server: &testServer{},
		Name:   "ServerName",
	}

	testStopServer(t, server)
}

func testStopServer(t *testing.T, server serverz.Server) {
	serverManager := serverz.NewManager()

	var wg sync.WaitGroup

	f := serverManager.StopServer(server, &wg)

	ctx := context.Background()
	f(ctx)

	var called bool

	if s, ok := server.(*serverz.NamedServer); ok {
		called = s.Server.(*testServer).ShutdownCalled
	} else {
		called = server.(*testServer).ShutdownCalled
	}

	if got, want := called, true; got != want {
		t.Fatal("serverz: Server.Shutdown should be called")
	}
}
