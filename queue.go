package serverz

import (
	"context"
	"sync"
)

type serverItem struct {
	server Server
	addr   string
}

// ServerQueue holds a list of servers and starts them at once.
type ServerQueue struct {
	servers []*serverItem

	manager Manager
}

// Append appends a new server to the list of existing ones.
func (q *ServerQueue) Append(server Server, addr string) {
	q.servers = append(
		q.servers,
		&serverItem{
			server: server,
			addr:   addr,
		},
	)
}

// Prepend prepends a new server to the list of existing ones.
func (q *ServerQueue) Prepend(server Server, addr string) {
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
func (q *ServerQueue) Start() <-chan error {
	ch := make(chan error, 2*len(q.servers))

	for _, server := range q.servers {
		go q.manager.ListenAndStartServer(server.server, server.addr)(ch)
	}

	return ch
}

// Stop tries to gracefully stop all the servers.
func (q *ServerQueue) Stop(ctx context.Context) {
	wg := &sync.WaitGroup{}

	for _, server := range q.servers {
		go q.manager.StopServer(server.server, wg)(ctx)
	}

	wg.Wait()
}
