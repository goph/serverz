package serverz

import (
	"context"
	"net"
	"sync"
)

// Queue holds a list of servers and starts them at once.
type Queue struct {
	servers []struct {
		server Server
		addr   net.Addr
	}
}

// NewQueue returns a new Queue.
func NewQueue() *Queue {
	return new(Queue)
}

// Append appends a new server to the list of existing ones.
func (q *Queue) Append(server Server, addr net.Addr) {
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
		if ls, ok := s.server.(listenServer); ok {
			go func(ch chan<- error, server listenServer, addr net.Addr) {
				ch <- server.ListenAndServe(addr)
			}(ch, ls, s.addr)
		} else {
			lis, err := listen(s.addr)
			if err != nil {
				ch <- err

				// Skip starting this server if we can't listen on the interface
				// Consuming to the error channel should end up being the application terminated anyway
				continue
			}

			go func(ch chan<- error, server Server, lis net.Listener) {
				ch <- server.Serve(lis)
			}(ch, s.server, lis)
		}
	}

	return ch
}

// Shutdown tries to gracefully stop all the servers.
func (q *Queue) Shutdown(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	merr := multiError{}

	for _, s := range q.servers {
		wg.Add(1)

		go func(server Server) {
			err := server.Shutdown(ctx)
			if err != nil {
				merr = append(merr, err)
			}

			wg.Done()
		}(s.server)
	}

	wg.Wait()

	return merr.ErrOrNil()
}

// Close immediately calls close for all servers.
func (q *Queue) Close() error {
	merr := multiError{}

	for _, s := range q.servers {
		err := s.server.Close()

		if err != nil {
			merr = append(merr, err)
		}
	}

	return merr.ErrOrNil()
}
