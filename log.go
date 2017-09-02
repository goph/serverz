package serverz

import (
	"fmt"
	"log"
	"strings"

	"github.com/goph/serverz/internal"
)

// DefaultLogger is used by servers as a default when no logger is specified.
var DefaultLogger logger = &defaultLogger{}

// DisableLogging disables the default logger.
func DisableLogging() {
	DefaultLogger = &nopLogger{}
}

// logger is the base interface for logging based on go-kit's logger.
type logger interface {
	Log(keyvals ...interface{}) error
}

type defaultLogger struct{}

// Log implements the logger interface.
func (*defaultLogger) Log(keyvals ...interface{}) error {
	ctx := internal.MapContext(keyvals)

	var line string

	for key, value := range ctx {
		line = fmt.Sprintf("%s %s=%q", line, key, value)
	}

	line = line + "\n"
	line = strings.TrimSpace(line)

	log.Print(line)

	return nil
}

type nopLogger struct{}

func (*nopLogger) Log(...interface{}) error { return nil }
