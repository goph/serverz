package named_test

import (
	"testing"

	"github.com/goph/serverz/mocks"
	"github.com/goph/serverz/named"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	s := &mocks.Server{}

	server := named.NewServer(s, "name")

	assert.Equal(t, s, server.Server)
	assert.Equal(t, "name", server.ServerName)
}

func TestServer_Name(t *testing.T) {
	server := &named.Server{
		ServerName: "name",
	}

	assert.Equal(t, "name", server.Name())
}