package serverz

import (
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func TestLoggerOption(t *testing.T) {
	logger := log.NewNopLogger()

	opts := newOptions(Logger(logger))

	assert.Equal(t, logger, opts.Logger())
}
