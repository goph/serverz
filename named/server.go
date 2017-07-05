package named

import "github.com/goph/serverz"

// Server can be used to add a name to a server which is useful eg. in logs.
type Server struct {
	serverz.Server

	ServerName string
}

// NewServer returns a new Server.
func NewServer(server serverz.Server, name string) *Server {
	return &Server{server, name}
}

// Name returns the server's name.
func (s *Server) Name() string {
	return s.ServerName
}
