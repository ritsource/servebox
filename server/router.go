package server

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	db "github.com/ritwik310/servebox/database"
)

// IndexHandler ...
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Hello there..\n To read files visit /read/:filename, and to download directly visit /download/:filename"))
}

// readFile does the file reading for DownloadHandler and ReadHandler
func readFile(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	fn := r.URL.Query().Get("file")     // Reading File Title from the query string (file)
	pw := r.URL.Query().Get("password") // Reading Password from the query string (passphrase)

	fmt.Println("fn", fn)
	fmt.Println("pw", pw)

	// Checking if filename and Password exist in th equery string
	if fn == "" || pw == "" {
		return []byte{}, errors.New("Both Filename and Password required in the query string, /read?file=example.txt&passowrd=examplepass")
	}

	file := db.File{Title: fn, Password: pw} // File Struct

	// Reading File
	fp, err := file.GetFile() // Returns file-path (fp)
	if err != nil {
		return []byte{}, err
	}

	// Reading the actual file
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

// ReadHandler ...
func ReadHandler(w http.ResponseWriter, r *http.Request) {
	b, err := readFile(w, r) // Read file
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	buf := bytes.NewBuffer(b) // New Buffer

	// Writing buf
	if _, err := buf.WriteTo(w); err != nil {
		fmt.Fprintf(w, "%s", err)
	}
}

// DownloadHandler ...
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	b, err := readFile(w, r) // Read file
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	buf := bytes.NewBuffer(b)                                  // New Buffer
	w.Header().Set("Content-type", "application/octet-stream") // Header Content-Type = "application/octet-stream"

	// Writing buf
	if _, err := buf.WriteTo(w); err != nil {
		fmt.Fprintf(w, "%s", err)
	}
}
