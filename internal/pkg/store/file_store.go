package store

import (
	"fmt"
	"io/ioutil"
	"os"
)

// FileStore represents store that saves data into a file.
// It implements Store inteface.
type FileStore struct {
	path string
}

// NewFileStore builds new FileStore.
func NewFileStore(path string) Store {
	return FileStore{
		path: path,
	}
}

// Save writes data into a file.
func (s FileStore) Save(data []byte) error {
	file, err := os.Create(s.path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("write to file: %w", err)
	}

	err = file.Sync()
	if err != nil {
		return fmt.Errorf("sync file: %w", err)
	}

	return nil
}

// Load loads data from a file.
func (s FileStore) Load() ([]byte, error) {
	file, err := os.Open(s.path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return data, nil
}
