package aio

import (
	"github.com/goph/serverz"
	"github.com/goph/serverz/internal/ext"
)

// Server implements all kinds of functionalities, such as:
//	- serverz.namer interface
//	- custom closers
type Server struct {
	serverz.Server

	Name   string
	Closer ext.Closer
}

// GetName returns the name of a server.
func (s *Server) GetName() string {
	return s.Name
}

// Close invokes the wrapped server's closer first then the ones from s.Closer if any.
func (s *Server) Close() error {
	closers := ext.Closers{s.Server}

	if s.Closer != nil {
		closers = append(closers, s.Closer)
	}

	return closers.Close()
}
