package grpc

import (
	"context"
	"sync"

	"google.golang.org/grpc"
)

// Server implements the Server interface for a gRPC server
type Server struct {
	*grpc.Server
}

// Shutdown implements Server interface
func (s *Server) Shutdown(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	c := make(chan struct{})

	go func() {
		defer close(c)
		s.GracefulStop()
		wg.Wait()
	}()

	select {
	case <-c:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close implements Server
func (s *Server) Close() error {
	s.Server.Stop()

	return nil
}
