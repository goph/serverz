package serverz

import (
	"context"
	"net"
)

// Server in this context is an abstraction over anything that can be started and stopped (either gracefully or forcefully)
// Typically they accept connections and serve over network, like HTTP or RPC servers
type Server interface {
	// Serve accepts incoming connections on the Listener l.
	//
	// Serve always returns a non-nil error.
	Serve(l net.Listener) error

	// Shutdown gracefully shuts down the server without interrupting any
	// active connections. Shutdown works by first closing all open
	// listeners, then closing all idle connections, and then waiting
	// indefinitely for connections to return to idle and then shut down.
	// If the provided context expires before the shutdown is complete,
	// then the context's error is returned.
	Shutdown(ctx context.Context) error

	// Close immediately closes all active net.Listeners.
	//
	// For a graceful shutdown, use Shutdown.
	Close() error
}
