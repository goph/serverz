package serverz

import (
	"net"
)

// AddrServer holds a net.Addr which carries the information of the listening network and address.
type AddrServer interface {
	Server

	Addr() net.Addr
}

type addr struct {
	network string
	addr    string
}

func (a *addr) Network() string {
	return a.network
}

func (a *addr) String() string {
	return a.addr
}

type addrServer struct {
	Server

	addr net.Addr
}

func (s *addrServer) Addr() net.Addr {
	return s.addr
}
