package serverz

import (
	"context"
	"net"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// AppServer wraps a server and provides all kinds of functionalities, such as:
//
//	- Logging server events
//	- Registering custom closers in the server
//	- Listening on a net.Addr
type AppServer struct {
	Server

	Name   string
	Addr   net.Addr
	Closer Closer

	Logger log.Logger
}

// Serve calls the underlying server.
func (s *AppServer) Serve(l net.Listener) error {
	level.Info(s.logger()).Log(
		"msg", "Starting server",
		"addr", l.Addr(),
		"server", s.Name,
	)

	return s.Server.Serve(l)
}

// Shutdown attempts to gracefully shut the underlying server down.
func (s *AppServer) Shutdown(ctx context.Context) error {
	level.Info(s.logger()).Log(
		"msg", "Attempting to shut server gracefully down",
		"server", s.Name,
	)

	return s.Server.Shutdown(ctx)
}

// Close invokes the wrapped server's closer first then the ones from s.Closer if any.
func (s *AppServer) Close() error {
	level.Info(s.logger()).Log(
		"msg", "Closing server",
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

	level.Info(s.logger()).Log(
		"msg", "Listening on address",
		"addr", addr,
		"server", s.Name,
	)

	return s.Serve(lis)
}

// logger returns the configured logger instance or the default logger.
func (s *AppServer) logger() log.Logger {
	if s.Logger == nil {
		return DefaultLogger
	}

	return s.Logger
}
