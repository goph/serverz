package serverz

// namer exposes a name for a server.
type namer interface {
	// GetName returns the name of a server.
	GetName() string
}

// getServerName extracts the name of the server if it implements namer.
func getServerName(server Server) (name string) {
	if server, ok := server.(namer); ok {
		name = server.GetName()
	}

	return
}
