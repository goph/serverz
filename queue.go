package serverz

import (
	"context"
	"net"
	"sync"

	"github.com/goph/serverz/internal/errors"
)

// Queue holds a list of servers and starts them at once.
type Queue struct {
	servers []struct {
		server Server
		addr   net.Addr
	}

	manager *Manager
}

// NewQueue returns a new Queue.
func NewQueue(opt ...Option) *Queue {
	return &Queue{
		manager: NewManager(opt...),
	}
}

// Append appends a new server to the list of existing ones.
func (q *Queue) Append(server Server, addr net.Addr) {
	// When no addr is provided try getting one from the server
	if s, ok := server.(AddrServer); ok && addr == nil {
		addr = s.GetAddr()
	}

	q.servers = append(
		q.servers,
		struct {
			server Server
			addr   net.Addr
		}{
			server,
			addr,
		},
	)
}

// Prepend prepends a new server to the list of existing ones.
func (q *Queue) Prepend(server Server, addr net.Addr) {
	// When no addr is provided try getting one from the server
	if s, ok := server.(AddrServer); ok && addr == nil {
		addr = s.GetAddr()
	}

	q.servers = append(
		[]struct {
			server Server
			addr   net.Addr
		}{
			{
				server,
				addr,
			},
		},
		q.servers...,
	)
}

// Start starts all the servers.
func (q *Queue) Start() <-chan error {
	ch := make(chan error, len(q.servers))

	for _, s := range q.servers {
		starter, err := q.manager.ListenAndStartServer(s.server, s.addr)
		if err != nil {
			panic(err)
		}

		go starter(ch)
	}

	return ch
}

// Shutdown tries to gracefully stop all the servers.
func (q *Queue) Shutdown(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	errBuilder := errors.NewMultiErrorBuilder()
	errBuilder.Message = "An error ocurred during server shutdown"

	for _, s := range q.servers {
		wg.Add(1)

		go func(server Server) {
			err := q.manager.StopServer(server, wg)(ctx)
			errBuilder.Add(err)

			wg.Done()
		}(s.server)
	}

	wg.Wait()

	return errBuilder.ErrOrNil()
}

// Close immediately calls close for all servers.
func (q *Queue) Close() error {
	errBuilder := errors.NewMultiErrorBuilder()
	errBuilder.Message = "An error ocurred during server close"

	for _, s := range q.servers {
		err := s.server.Close()
		errBuilder.Add(err)
	}

	return errBuilder.ErrOrNil()
}
