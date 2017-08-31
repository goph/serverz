package serverz

import (
	"errors"
	"net"
)

// AddrServer holds a net.Addr which carries the information of the listening network and address.
type AddrServer interface {
	Server

	GetAddr() net.Addr
}

// virtualAddr is a fake, in-memory address representation.
type virtualAddr struct {
	network string
	addr    string
}

// NewAddr returns a new in-memory Addr.
func NewAddr(network, addr string) net.Addr {
	return &virtualAddr{network, addr}
}

// Network returns the address's network name.
func (a *virtualAddr) Network() string {
	return a.network
}

// String returns the address's string representation.
func (a *virtualAddr) String() string {
	return a.addr
}

// virtualListener is a fake, in-memory listener.
// It's meant to be used with "servers" with no network activity involved (like daemons, etc).
type virtualListener struct {
	addr net.Addr
}

// newVirtualListener returns a fake, in-memory listener.
func newVirtualListener(addr net.Addr) *virtualListener {
	return &virtualListener{addr}
}

// Accept immediately returns with an error since it's noop.
func (l *virtualListener) Accept() (net.Conn, error) {
	return nil, errors.New("This listener is not capable of serving any connections")
}

// Close is noop.
func (l *virtualListener) Close() error {
	return nil
}

// Addr returns the listener's network address.
func (l *virtualListener) Addr() net.Addr {
	return l.addr
}
