package named_test

import (
	"testing"

	"github.com/goph/serverz"
	"github.com/goph/serverz/internal"
	"github.com/goph/serverz/internal/mocks"
	"github.com/goph/serverz/named"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	s := &mocks.Server{}

	server := named.NewServer(s, "name")

	assert.Equal(t, s, server.Server)
	assert.Equal(t, "name", server.Name)
}

func TestServerIsAServer(t *testing.T) {
	assert.Implements(t, (*serverz.Server)(nil), new(named.Server))
}

func TestServer_Name(t *testing.T) {
	server := &named.Server{
		Name: "name",
	}

	assert.Implements(t, (*internal.NamedServer)(nil), server)
	assert.Equal(t, "name", server.GetName())
}
