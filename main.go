package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	db "github.com/ritwik310/servebox/database"
)

func init() {
	// Creates required directories, doesn't delete if already exist
	err := os.MkdirAll(db.FileLoc, os.ModePerm)
	err = os.MkdirAll(db.PassLoc, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

// HelloServer ...
func HelloServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

// FileServer ...
func FileServer(w http.ResponseWriter, r *http.Request) {
	pp := r.URL.Query().Get("passphrase")
	fmt.Println("passphrase", pp)

	f := db.File{
		Location: "index2.txt",
		Password: pp,
	}

	p, err := f.GetFile()
	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	} else if p == "" {
		fmt.Fprintf(w, "%s", "Wrong Passphrase!")
		return
	}

	data, err := ioutil.ReadFile(p)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	buf := bytes.NewBuffer(data)

	// w.Header().Set("Content-type", "application/octet-stream")

	if _, err := buf.WriteTo(w); err != nil {
		fmt.Fprintf(w, "%s", err)
	}
}

func main() {
	fmt.Println(db.BaseLoc)
	http.HandleFunc("/hello", HelloServer)
	http.HandleFunc("/", FileServer)

	// err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
