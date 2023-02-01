package main

import (
	"sync"

	"github.com/google/uuid"
)

// Cluster describes
type Cluster struct {
	servers map[string]*Storage // map[serverAddress]*Storage
	mu      sync.RWMutex
}

// NewCluster creates a new Cluster
func NewCluster() *Cluster {
	return &Cluster{
		servers: make(map[string]*Storage),
	}
}

// AddServer adds a server to the cluster
func (c *Cluster) AddServer(serverAddress string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.servers[serverAddress] = NewStorage()
}

// GetServer returns a server from the cluster
func (c *Cluster) GetServer(serverAddress string) (*Storage, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	server, ok := c.servers[serverAddress]
	return server, ok
}

// DeleteServer deletes a server from the cluster
func (c *Cluster) DeleteServer(serverAddress string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.servers, serverAddress)
}

// GetChunk returns a chunk from a server in the cluster
func (c *Cluster) GetChunk(serverAddress string, key uuid.UUID) ([]byte, bool) {
	server, ok := c.GetServer(serverAddress)
	if !ok {
		return nil, false
	}
	return server.Get(key)
}

// SetChunk sets a chunk in a server in the cluster
func (c *Cluster) SetChunk(serverAddress string, key uuid.UUID, chunk []byte) {
	server, ok := c.GetServer(serverAddress)
	if !ok {
		return
	}
	server.Set(key, chunk)
}
