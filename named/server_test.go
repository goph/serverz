package named_test

import (
	"testing"

	"github.com/goph/serverz/mocks"
	"github.com/goph/serverz/named"
)

func TestNewServer(t *testing.T) {
	s := &mocks.Server{}

	server := named.NewServer(s, "name")

	if server.Server != s {
		t.Error("expected the Server member to be set to variable 's'")
	}

	if got, want := server.ServerName, "name"; got != want {
		t.Errorf("expected the server's name to be 'name', received: '%s'", server.ServerName)
	}
}
