package serverz

import (
	"net"
)

// AddrServer holds a net.Addr which carries the information of the listening network and address.
type AddrServer interface {
	Server

	GetAddr() net.Addr
}
