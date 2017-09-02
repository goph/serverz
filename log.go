package serverz

import (
	"os"

	"github.com/go-kit/kit/log"
)

// DefaultLogger is used by servers as a default when no logger is specified.
var DefaultLogger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))

// DisableLogging disables logging for servers.
func DisableLogging() {
	DefaultLogger = log.NewNopLogger()
}
