package server

import (
	"net/http"
)

// Start starts up the public server
func Start(port string) error {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/read", ReadHandler)
	http.HandleFunc("/download", DownloadHandler)

	// err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)
	return http.ListenAndServe(":"+port, nil)
}
