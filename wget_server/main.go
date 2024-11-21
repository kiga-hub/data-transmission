package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// run
// ./test -fileDir=/home/workspace/golang/data-transmission/wget_server -port=4567

// test
// wget http://192.168.8.244:4567/download/data.tar.gz -O source/data.tar.gz

func main() {
	fileDir := flag.String("fileDir", "/home/workspace/golang/data-transmission/wget_server", "Directory to serve files from")
	port := flag.String("port", "4567", "Port to listen on")
	flag.Parse()

	// handle download
	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		fileName := strings.TrimPrefix(r.URL.Path, "/download/")
		filePath := filepath.Join(*fileDir, fileName)

		// set header
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", "application/octet-stream")

		// serve file
		http.ServeFile(w, r, filePath)
	})

	// start server
	log.Printf("Start Server, Listen port: %s", *port)
	if err := http.ListenAndServe("0.0.0.0:"+*port, nil); err != nil {
		log.Fatalf("Start failed: %v", err)
	}
}
