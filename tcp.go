package serverz

// NewTCPServer returns a server which contains address information.
func NewTCPServer(s Server, a string) AddrServer {
	return &addrServer{
		Server: s,
		addr: &addr{
			network: "tcp",
			addr:    a,
		},
	}
}
