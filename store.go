package main

import (
	"sync"

	"github.com/google/uuid"
)

// Storage describes temporary storage for chunks
type Storage struct {
	db map[uuid.UUID][]byte // map[key]chunk
	mu sync.RWMutex
}

// NewStorage creates a new Storage
func NewStorage() *Storage {
	return &Storage{
		db: make(map[uuid.UUID][]byte),
	}
}

// Get returns a chunk from the storage
func (s *Storage) Get(key uuid.UUID) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	chunk, ok := s.db[key]
	return chunk, ok
}

// Set sets a chunk in the storage
func (s *Storage) Set(key uuid.UUID, chunk []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.db[key] = chunk
}

// Delete deletes a chunk from the storage
func (s *Storage) Delete(key uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.db, key)
}
