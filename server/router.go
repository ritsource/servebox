package server

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	db "github.com/ritwik310/servebox/database"
)

const html = `
<!DOCTYPE html>
<head>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>ServeBox</title>
</head>
<body>
	<form>
		<h3>ServeBox</h3>
    <label>File:</label>
    <input id="file-input" name="file"/>
    <br/>
    <label>Password:</label>
    <input id="password-input" name="password"/>
    <br/>

    <button type="button" id="submit-button" >Submit</button>
    <button type="button" id="download-button" >Download</button>
  </form>

  <script>
		document.getElementById('submit-button').onclick = function() {
			const file = document.getElementById('file-input').value
			const password = document.getElementById('password-input').value

			window.location.href = '/read?file='+encodeURIComponent(file)+"&password="+password
			return false;
		};

		document.getElementById('download-button').onclick = function() {
			const file = document.getElementById('file-input').value
			const password = document.getElementById('password-input').value
			
			var win = window.open('/download?file='+encodeURIComponent(file)+"&password="+password, '_blank');
			win.focus();
			
			return false;
		};

  </script>
</body>
</html>`

// IndexHandler ...
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
	// w.Write([]byte("Hello there..\n To read files visit /read/:filename, and to download directly visit /download/:filename"))
}

// readFile does the file reading for DownloadHandler and ReadHandler
func readFile(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	fn, err := url.QueryUnescape(r.URL.Query().Get("file")) // Reading File Title from the query string (file)
	pw := r.URL.Query().Get("password")                     // Reading Password from the query string (passphrase)
	if err != nil {
		return []byte{}, err
	}

	fmt.Println(fn)
	fmt.Println(pw)

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

	if err != nil && err.Error() == "wrong:password" {
		fmt.Fprintf(w, "%s", "Wrong Password!")
		return
	} else if err != nil {
		fmt.Fprintf(w, "%s", err)
		// fmt.Fprintf(w, "%s", "File cannot be found, or maybe internal server error")
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

	if err != nil && err.Error() == "wrong:password" {
		fmt.Fprintf(w, "%s", "Wrong Password!")
		return
	} else if err != nil {
		fmt.Fprintf(w, "%s", err)
		// fmt.Fprintf(w, "%s", "File cannot be found, or maybe internal server error")
		return
	}

	buf := bytes.NewBuffer(b)                                  // New Buffer
	w.Header().Set("Content-type", "application/octet-stream") // Header Content-Type = "application/octet-stream"

	// Writing buf
	if _, err := buf.WriteTo(w); err != nil {
		fmt.Fprintf(w, "%s", err)
	}
}
