package serverz

import (
	"context"
	"sync"

	"google.golang.org/grpc"
)

// GrpcServer implements the Server interface for a gRPC server
type GrpcServer struct {
	*grpc.Server
}

// Shutdown implements Server interface
func (s *GrpcServer) Shutdown(ctx context.Context) error {
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
func (s *GrpcServer) Close() error {
	s.Server.Stop()

	return nil
}
