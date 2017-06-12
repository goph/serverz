package serverz

import (
	"context"
	"testing"
)

func TestNamedServer_Serve(t *testing.T) {
	spy := &testServer{}
	server := &NamedServer{
		Server: spy,
	}

	server.Serve(&testListener{})

	if spy.ServeCalled != true {
		t.Fatal("Server.Serve should be called")
	}
}

func TestNamedServer_Shutdown(t *testing.T) {
	spy := &testServer{}
	server := &NamedServer{
		Server: spy,
	}

	server.Shutdown(context.Background())

	if spy.ShutdownCalled != true {
		t.Fatal("Server.Shutdown should be called")
	}
}

func TestNamedServer_Close(t *testing.T) {
	spy := &testServer{}
	server := &NamedServer{
		Server: spy,
	}

	server.Close()

	if spy.CloseCalled != true {
		t.Fatal("Server.Close should be called")
	}
}
