package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// Upload describes upload handler
type Upload struct {
	db      *DB
	cluster *Cluster
	servers []string
}

// NewUpload creates a new Upload
func NewUpload(db *DB, cluster *Cluster, servers []string) *Upload {
	return &Upload{
		db:      db,
		cluster: cluster,
		servers: servers,
	}
}

func (h *Upload) Handler(w http.ResponseWriter, r *http.Request) {
	// function body of a http.HandlerFunc
	r.Body = http.MaxBytesReader(w, r.Body, 100<<20+1024) // Max 100MB + 1KB
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// parse file field
	p, err := reader.NextPart()
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if p.FormName() != "file" {
		http.Error(w, "file is expected", http.StatusBadRequest)
		return
	}

	var (
		eg     errgroup.Group
		stop   bool
		stopMu sync.Mutex
		chunks []uuid.UUID // map[serverAddress]chunk_keys
	)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	for i := 0; ; i++ {
		if i == len(h.servers) {
			i = 0
		}

		if stop {
			cancel()
			break
		}

		// Read chunkSize bytes from the file
		buf := make([]byte, chunkSize)
		n, err := p.Read(buf)
		if err != nil && err != io.EOF {
			cancel()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		if n == 0 {
			break
		}

		if stop {
			cancel()
			break
		}

		// recreate variables to prevent rewriting
		i := i
		chunkUUID := uuid.New()
		chunks = append(chunks, chunkUUID)

		eg.Go(func() error {
			err := h.processChunk(ctx, buf[:n], chunkUUID, i)
			if err != nil {
				stopMu.Lock()
				stop = true
				stopMu.Unlock()
			}
			return err
		})
	}

	err = eg.Wait()
	if err != nil {
		log.Printf("Error processing chunk: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.db.SetFile(p.FileName(), &File{
		fileName:    p.FileName(),
		contentType: p.Header.Get("Content-Type"),
		chunks:      chunks,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *Upload) processChunk(ctx context.Context, chunk []byte, chunkUUID uuid.UUID, serverIndex int) error {
	// we can use context to cancel the process

	// get the server address
	serverAddress := h.servers[serverIndex]

	// save relationship between chunk and server in the database
	h.db.SetChunkServer(chunkUUID, serverAddress)

	// upload the chunk to the storage server
	h.cluster.SetChunk(serverAddress, chunkUUID, chunk)

	return nil
}
