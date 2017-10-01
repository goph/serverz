package serverz

import (
	"context"
	"net"
)

// AppServer wraps a server and provides all kinds of functionalities, such as:
//
//	- Logging server events
//	- Registering custom closers in the server
//	- Listening on a net.Addr
type AppServer struct {
	Server

	// Name is used in logs to identify the server.
	Name string

	// Addr is optionally used for listening if specified.
	Addr net.Addr

	// Closer is called together with the servers regular Close method.
	Closer Closer

	// Logger specifies an optional logger.
	// If nil, logging goes to os.Stderr via the log package's
	// standard logger.
	Logger logger
}

// Serve calls the underlying server.
func (s *AppServer) Serve(l net.Listener) error {
	s.logger().Log(
		"msg", "starting server",
		"addr", l.Addr(),
		"server", s.Name,
	)

	return s.Server.Serve(l)
}

// Shutdown attempts to gracefully shut the underlying server down.
func (s *AppServer) Shutdown(ctx context.Context) error {
	s.logger().Log(
		"msg", "shutting server gracefully down",
		"server", s.Name,
	)

	return s.Server.Shutdown(ctx)
}

// Close invokes the wrapped server's closer first then the ones from s.Closer if any.
func (s *AppServer) Close() error {
	s.logger().Log(
		"msg", "closing all remaining connections",
		"server", s.Name,
	)

	closers := Closers{s.Server}

	if s.Closer != nil {
		closers = append(closers, s.Closer)
	}

	return closers.Close()
}

// ListenAndServe listens on a network address and then
// calls Serve to handle requests on incoming connections.
func (s *AppServer) ListenAndServe(addr net.Addr) error {
	// If no address is passed, try using the one configured in the server
	if addr == nil {
		addr = s.Addr
	}

	lis, err := listen(addr)
	if err != nil {
		return err
	}

	s.logger().Log(
		"msg", "listening on address",
		"addr", lis.Addr().String(),
		"server", s.Name,
	)

	return s.Serve(lis)
}

// logger returns the configured logger instance or the default logger.
func (s *AppServer) logger() logger {
	if s.Logger == nil {
		return DefaultLogger
	}

	return s.Logger
}
