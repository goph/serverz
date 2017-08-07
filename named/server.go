package named

import "github.com/goph/serverz"

// Server can be used to add a name to a server which is useful eg. in logs.
type Server struct {
	serverz.Server

	Name string
}

// NewServer returns a new Server.
func NewServer(server serverz.Server, name string) *Server {
	return &Server{server, name}
}

// GetName returns the name of a server.
func (s *Server) GetName() string {
	return s.Name
}
