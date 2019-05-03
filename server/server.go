package server

import (
	"log"
	"net/http"
)

// Port is the port that server will listen to
var Port = 6060

// Start starts up the public server
func Start() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/read", ReadHandler)
	http.HandleFunc("/download", DownloadHandler)

	// err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)
	log.Fatal(http.ListenAndServe(":6060", nil))
}
