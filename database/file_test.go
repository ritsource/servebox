package database_test

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"

	db "github.com/ritwik310/servebox/database"
)

var SourceDir string

// HandleError1 handles errors by exiting the process
func HandleError1(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Creates required directories, doesn't delete if already exist
	err := os.MkdirAll(db.FileLoc, os.ModePerm)
	err = os.MkdirAll(db.PassLoc, os.ModePerm)
	HandleError1(err)

	SourceDir = path.Join("..", "test_source")

	ClearSB()
	PopSrc()
}

// ClearSB clears a out unnecessary files from db.Fileloc
func ClearSB() {
	// Open the FileLocation Dir
	dir, err := os.Open(db.FileLoc)
	HandleError1(err)
	defer dir.Close()

	// Sub-Directories
	names, err := dir.Readdirnames(-1)
	HandleError1(err)

	// Deleting Sub-Directories
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(db.FileLoc, name))
		HandleError1(err)
	}

}

// PopSrc populates necessary files in the "test_source" dir
func PopSrc() {
	err := os.MkdirAll(SourceDir, os.ModePerm)
	HandleError1(err)

	files := map[string]string{
		"test_doc.txt": "test_doc",
	}

	for k, v := range files {
		err = ioutil.WriteFile(path.Join(SourceDir, k), []byte(v), 0644)
		HandleError1(err)
	}
}

func TestCopyFile(t *testing.T) {
	src := path.Join(SourceDir, "test_doc.txt") // Source File Path

	// File Struct, containing Title and password
	file := db.File{
		Title:    "test_doc.txt",
		Password: "mypassword1",
	}

	// NOTE: Always run copyFile before copyFileDup
	copyFile(t, file, src)    // Copy file
	copyFileDup(t, file, src) // Copy same file again
	copyFileDup(t, file, src) // Copy again, but this will be deleted
	copyFileRename(t, file, src, "new_doc.txt")

	// File Content Map, NOTE: "test_doc.copy.copy.txt" has been Deleted
	fcmap := map[string]string{
		"test_doc.txt":      "test_doc",
		"test_doc.copy.txt": "test_doc",
		"new_doc.txt":       "test_doc",
	}

	checkCreatedFiles(t, "mypassword1", fcmap) // Checking File Content

	// Testing File Deletion
	f2bDeleted := db.File{Title: "test_doc.copy.copy.txt", Password: "mypassword1"} // File to be deleted
	removeFile(t, f2bDeleted)                                                       // Deleting f2bDeleted ("test_doc.copy.copy.txt")

	// Check if deleted or not
	_, err := ioutil.ReadFile(path.Join(db.FileLoc, f2bDeleted.Password, f2bDeleted.Title))
	if err == nil {
		t.Error("Error:", "Unable to Delete file")
	}
}

func copyFile(t *testing.T, file db.File, src string) {
	_, err := file.CopyFile(src)
	if err != nil && err.Error() != "dup:err" {
		t.Error("Error:", err)
	}
}

func copyFileDup(t *testing.T, file db.File, src string) {
	_, err := file.CopyFileDup(src)
	if err != nil {
		t.Error("Error:", err)
	}
}

func copyFileRename(t *testing.T, file db.File, src string, newname string) {
	_, err := file.CopyFileRename(src, newname)
	if err != nil && err.Error() != "dup:err" {
		t.Error("Error:", err)
	}
}

func removeFile(t *testing.T, file db.File) {
	err := file.RemoveFile()
	if err != nil {
		t.Error("Error:", err)
	}
}

func checkCreatedFiles(t *testing.T, password string, fcmap map[string]string) {
	for k, v := range fcmap {
		// filesl
		b, err := ioutil.ReadFile(path.Join(db.FileLoc, password, k))
		if err != nil {
			t.Error("Error:", err)
		}

		if string(b) != v {
			t.Error("Error:", errors.New("Content Mismatch, "+path.Join(db.FileLoc, password, k)))
		}
	}
}
