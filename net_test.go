package serverz_test

import (
	"testing"

	. "github.com/goph/serverz"
	"github.com/stretchr/testify/assert"
)

func TestNewAddr(t *testing.T) {
	addr := NewAddr("network", "addr")

	assert.Equal(t, "network", addr.Network())
	assert.Equal(t, "addr", addr.String())
}
