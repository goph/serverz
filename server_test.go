package serverz_test

import (
	"context"
	"testing"

	"github.com/goph/serverz"
)

func TestNamedServerServe(t *testing.T) {
	spy := &TestServer{}
	ns := &serverz.NamedServer{
		Server: spy,
	}

	ns.Serve(&TestListener{})

	if got, want := spy.ServeCalled, true; got != want {
		t.Fatal("server: Server.Serve should be called")
	}
}

func TestNamedServerShutdown(t *testing.T) {
	spy := &TestServer{}
	ns := &serverz.NamedServer{
		Server: spy,
	}

	ns.Shutdown(context.Background())

	if got, want := spy.ShutdownCalled, true; got != want {
		t.Fatal("server: Server.Shutdown should be called")
	}
}

func TestNamedServerClose(t *testing.T) {
	spy := &TestServer{}
	ns := &serverz.NamedServer{
		Server: spy,
	}

	ns.Close()

	if got, want := spy.CloseCalled, true; got != want {
		t.Fatal("server: Server.Close should be called")
	}
}
