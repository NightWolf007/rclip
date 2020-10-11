// Package clipboard provides system clipboard management methods.
package clipboard

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
)

// Clipboard represents the basic clipboard interface.
// It can be used primary for mocking.
type Clipboard interface {
	Read() ([]byte, error)
	Write([]byte) error
}

type clipboardImpl struct{}

func (c clipboardImpl) Read() ([]byte, error) {
	val, err := clipboard.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("clipboard ReadAll: %w", err)
	}

	return []byte(val), nil
}

func (c clipboardImpl) Write(val []byte) error {
	err := clipboard.WriteAll(string(val))
	if err != nil {
		return fmt.Errorf("clipboard WriteAll: %w", err)
	}

	return nil
}

// Watch creates a watch stream for clipboard changes.
func Watch() *Stream {
	return &Stream{
		RequestDelay: time.Second,
		Clipboard:    clipboardImpl{},
		TimeAfter:    time.After,
	}
}

// Read reads system clipboard.
func Read() ([]byte, error) {
	return clipboardImpl{}.Read()
}

// Write writes new value into system clipboard.
func Write(val []byte) error {
	return clipboardImpl{}.Write(val)
}

// IsSupported returns truw if system clipboard unsupoorted.
func IsSupported() bool {
	return !clipboard.Unsupported
}
