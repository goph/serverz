package serverz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerOption(t *testing.T) {
	logger := &defaultLogger{}

	opts := newOptions(Logger(logger))

	assert.Equal(t, logger, opts.Logger())
}
