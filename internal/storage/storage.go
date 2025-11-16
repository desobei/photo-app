// Package storage implements the Adapter pattern.
package storage

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

// Storage defines persistence operations.
type Storage interface {
	Save(id string, data []byte) error
	Load(id string) ([]byte, error)
}

// MapAdapter adapts a map to the Storage interface.
type MapAdapter struct {
	data map[string][]byte
}

// NewMapAdapter creates a new map-based storage adapter.
func NewMapAdapter() *MapAdapter {
	return &MapAdapter{
		data: make(map[string][]byte),
	}
}

// Save stores data in the map.
func (m *MapAdapter) Save(id string, data []byte) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if data == nil {
		return fmt.Errorf("data cannot be nil")
	}
	m.data[id] = data
	return nil
}

// Load retrieves data from the map.
func (m *MapAdapter) Load(id string) ([]byte, error) {
	data, ok := m.data[id]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, id)
	}
	return data, nil
}
