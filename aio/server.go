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

	ServerName string
	Closers    ext.Closers
}

// Name returns the server's name.
func (s *Server) Name() string {
	return s.ServerName
}

// Close invokes the wrapped server's closer first then the ones from s.Closers if any.
func (s *Server) Close() error {
	closers := append(ext.Closers{s.Server}, s.Closers...)

	return closers.Close()
}
