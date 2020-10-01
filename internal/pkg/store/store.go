// Package store provides store implementation.
package store

// Store represents an abstract byte storage.
type Store interface {
	Save([]byte) error
	Load() ([]byte, error)
}
