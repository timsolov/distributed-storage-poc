package main

import (
	"sync"

	"github.com/google/uuid"
)

// File describes files table in database
// +------------------+------------------+------------------+
// | file_name        | content_type     | chunks[]         |
// +------------------+------------------+------------------+
// | file1.jpg        | image/jpeg       | UUID1, UUID2     |
// | file2.jpeg       | image/jpeg       | UUID3            |
// +------------------+------------------+------------------+
type File struct {
	fileName    string
	contentType string
	chunks      []uuid.UUID
}

// ChunksServers describes chunks_servers table in database
// it is a map of chunk key to server address
// +------------------+------------------+
// | chunk_key        | server_address   |
// +------------------+------------------+
// | UUID1            | server1          |
// | UUID2            | server2          |
// | UUID3            | server3          |
// +------------------+------------------+
type ChunksServers map[uuid.UUID]string // map[chunkKey]serverAddress

// DB describes a database
type DB struct {
	files         map[string]*File // map[fileName]File
	chunksServers ChunksServers    // map[chunkKey]serverAddress
	mu            sync.RWMutex
}

// NewDB creates a new DB
func NewDB() *DB {
	return &DB{
		files:         make(map[string]*File),
		chunksServers: make(ChunksServers),
	}
}

// GetFile returns a file from the database
func (db *DB) GetFile(fileName string) (*File, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	file, ok := db.files[fileName]
	return file, ok
}

// SetFile sets a file in the database
func (db *DB) SetFile(fileName string, file *File) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.files[fileName] = file
}

// DeleteFile deletes a file from the database
func (db *DB) DeleteFile(fileName string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.files, fileName)
}

// GetChunkServer returns a server address for a chunk
func (db *DB) GetChunkServer(chunkKey uuid.UUID) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	serverAddress, ok := db.chunksServers[chunkKey]
	return serverAddress, ok
}

// SetChunkServer sets a server address for a chunk
func (db *DB) SetChunkServer(chunkKey uuid.UUID, serverAddress string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.chunksServers[chunkKey] = serverAddress
}
