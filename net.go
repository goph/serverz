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

// addr is a fake, in-memory address representation.
type addr struct {
	network string
	addr    string
}

// NewAddr returns a new in-memory Addr.
func NewAddr(network, address string) net.Addr {
	return addr{network, address}
}

// Network returns the address's network name.
func (a addr) Network() string {
	return a.network
}

// String returns the address's string representation.
func (a addr) String() string {
	return a.addr
}

// listener is a fake, in-memory listener.
// It's meant to be used with "servers" with no network activity involved (like daemons, etc).
type listener struct {
	addr net.Addr
}

// listen either listens on a real network interface
// or returns a fake, in-memory listener when no address is provided.
func listen(addr net.Addr) (net.Listener, error) {
	// Servers without an address can still work without network.
	if addr != nil {
		return net.Listen(addr.Network(), addr.String())
	}

	addr = NewAddr("none", "none")

	return &listener{addr}, nil
}

// NewListener returns a new Listener.
func NewListener(addr net.Addr) net.Listener {
	return &listener{addr}
}

// Accept immediately returns with an error since it's noop.
func (l *listener) Accept() (net.Conn, error) {
	return nil, errors.New("This listener is not capable of serving any connections")
}

// Close is noop.
func (l *listener) Close() error {
	return nil
}

// Addr returns the listener's network address.
func (l *listener) Addr() net.Addr {
	return l.addr
}
