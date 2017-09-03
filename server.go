package serverz

import (
	"context"
	"net"
)

// Server in this context is an abstraction over anything that can be started and stopped (either gracefully or forcefully).
// Typically they accept connections and serve over network, like HTTP servers.
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

// listenServer can listen on an interface on its own when provided with an address.
type listenServer interface {
	// ListenAndServe listens on a network address and then
	// calls Serve to handle requests on incoming connections.
	ListenAndServe(addr net.Addr) error
}
