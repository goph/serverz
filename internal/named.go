package internal

// NamedServer exposes a name for a server.
type NamedServer interface {
	// GetName returns the name of a server.
	GetName() string
}
