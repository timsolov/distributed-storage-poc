package main

import (
	"net/http"

	gmb "github.com/timsolov/gorilla-mux-booster"
)

// our storage servers
var storageServers = []string{"storage0", "storage1", "storage2", "storage3", "storage4"}

// chunkSize is the max size of each chunk
const chunkSize = 512 * 1024

func main() {
	cluster := NewCluster()
	db := NewDB()

	for i := 0; i < len(storageServers); i++ {
		cluster.AddServer(storageServers[i])
	}

	r := gmb.NewRouter()
	r.PUT("/upload", NewUpload(db, cluster, storageServers).Handler)
	r.GET("/download/{filename}", NewDownload(db, cluster).Handler)

	http.ListenAndServe(":8080", r)
}
