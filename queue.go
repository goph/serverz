package serverz

import (
	"context"
	"sync"

	"github.com/goph/serverz/internal/errors"
)

// Queue holds a list of servers and starts them at once.
type Queue struct {
	servers []AddrServer

	manager *Manager
}

// NewQueue returns a new Queue.
func NewQueue(opt ...Option) *Queue {
	return &Queue{
		manager: NewManager(opt...),
	}
}

// Append appends a new server to the list of existing ones.
func (q *Queue) Append(server AddrServer) {
	q.servers = append(q.servers, server)
}

// Prepend prepends a new server to the list of existing ones.
func (q *Queue) Prepend(server AddrServer) {
	q.servers = append(
		[]AddrServer{server},
		q.servers...,
	)
}

// Start starts all the servers.
func (q *Queue) Start() <-chan error {
	ch := make(chan error, len(q.servers))

	for _, server := range q.servers {
		starter, err := q.manager.ListenAndStartServer(server, server.GetAddr())
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

	for _, server := range q.servers {
		wg.Add(1)

		go func(server AddrServer) {
			err := q.manager.StopServer(server, wg)(ctx)
			errBuilder.Add(err)

			wg.Done()
		}(server)
	}

	wg.Wait()

	return errBuilder.ErrOrNil()
}

// Close immediately calls close for all servers.
func (q *Queue) Close() error {
	errBuilder := errors.NewMultiErrorBuilder()
	errBuilder.Message = "An error ocurred during server close"

	for _, server := range q.servers {
		err := server.Close()
		errBuilder.Add(err)
	}

	return errBuilder.ErrOrNil()
}
