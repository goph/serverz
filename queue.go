package serverz

import (
	"context"
	"sync"

	"github.com/goph/serverz/internal/errors"
)

type serverItem struct {
	server Server
	addr   string
}

// Queue holds a list of servers and starts them at once.
type Queue struct {
	servers []*serverItem

	manager *Manager
}

// NewQueue returns a new Queue.
func NewQueue(manager *Manager) *Queue {
	return &Queue{
		manager: manager,
	}
}

// Append appends a new server to the list of existing ones.
func (q *Queue) Append(server Server, addr string) {
	q.servers = append(
		q.servers,
		&serverItem{
			server: server,
			addr:   addr,
		},
	)
}

// Prepend prepends a new server to the list of existing ones.
func (q *Queue) Prepend(server Server, addr string) {
	q.servers = append(
		[]*serverItem{
			&serverItem{
				server: server,
				addr:   addr,
			},
		},
		q.servers...,
	)
}

// Start starts all the servers.
func (q *Queue) Start() <-chan error {
	ch := make(chan error, 2*len(q.servers))

	for _, server := range q.servers {
		go q.manager.ListenAndStartServer(server.server, "tcp", server.addr)(ch)
	}

	return ch
}

// Stop tries to gracefully stop all the servers.
func (q *Queue) Stop(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	errBuilder := errors.NewMultiErrorBuilder()
	errBuilder.Message = "An error ocurred during server shutdown"

	for _, server := range q.servers {
		wg.Add(1)

		go func(server *serverItem) {
			err := q.manager.StopServer(server.server, wg)(ctx)
			errBuilder.Add(err)

			wg.Done()
		}(server)
	}

	wg.Wait()

	return errBuilder.ErrOrNil()
}
