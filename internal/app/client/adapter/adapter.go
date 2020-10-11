// Package adapter provides server adapter implementation.
package adapter

// Adapter is an abstract RClip server adapter.
type Adapter interface {
	// Push sends data to the RClip server.
	Push(data []byte) error
	// Get returns the latest clipboard data.
	Get() ([]byte, error)
	// GetAll returns all clipboard history available.
	GetAll() ([][]byte, error)
}
