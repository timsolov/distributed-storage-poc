package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Download describes a download request handler
type Download struct {
	db      *DB
	cluster *Cluster
}

// NewDownload creates a new Download
func NewDownload(db *DB, cluster *Cluster) *Download {
	return &Download{
		db:      db,
		cluster: cluster,
	}
}

// http handler which joins chunks and returns the file
func (h *Download) Handler(w http.ResponseWriter, r *http.Request) {
	fileName := mux.Vars(r)["filename"]
	if fileName == "" {
		http.Error(w, "file is expected", http.StatusBadRequest)
		return
	}

	file, ok := h.db.GetFile(fileName)
	if !ok {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.fileName))

	for _, chunk := range file.chunks {
		serverAddress, ok := h.db.GetChunkServer(chunk)
		if !ok {
			http.Error(w, "chunk not found in database", http.StatusNotFound)
			return
		}

		b, ok := h.cluster.GetChunk(serverAddress, chunk)
		if !ok {
			http.Error(w, "chunk not found in storage", http.StatusNotFound)
			return
		}

		_, err := w.Write(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
