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

func TestGetFile(t *testing.T) {
	filename := "get_test_doc.txt"              // Name of file (Title)
	password := "mypassword1"                   // Password for File
	src := path.Join(SourceDir, "test_doc.txt") // Source File Path

	f := db.File{Title: filename} // File Struct, containing Title and password

	// Copying File from source
	np, _ := copyFile(t, f, src) // copyFile function takes care of error

	_, npFile := path.Split(np)                                            // Geting directory-name and file-name from np (new filepath)
	pw := db.Password{Title: npFile, Password: password, FileName: npFile} // New Password Struct

	// Writing the Password
	err := pw.Write()
	if err != nil {
		t.Error(err)
	}

	qFile := db.File{Title: filename, Password: password} // File Struct for Query

	fp, err := qFile.GetFile() // Reading file, using GetFile
	if err != nil {
		t.Error(err)
	}

	if fp != path.Join(db.FileLoc, filename) {
		t.Error(errors.New("Query from GetFile, doesn't match the expected path: " + fp + " != " + path.Join(db.FileLoc, filename)))
	}
}

func TestCopyAndRemove(t *testing.T) {
	src := path.Join(SourceDir, "test_doc.txt") // Source File Path

	// File Struct, containing Title and password
	f := db.File{
		Title:    "test_doc.txt",
		Password: "mypassword1",
	}

	// NOTE: Always run copyFile before copyFileDup
	copyFile(t, f, src)    // Copy file
	copyFileDup(t, f, src) // Copy same file again
	copyFileDup(t, f, src) // Copy again, but this will be deleted
	copyFileRename(t, f, src, "test_doc.new.txt")

	// File Content Map, NOTE: "test_doc.copy.copy.txt" has been Deleted
	fcmap := map[string]string{
		"test_doc.txt":      "test_doc",
		"test_doc.copy.txt": "test_doc",
		"test_doc.new.txt":  "test_doc",
	}

	checkCreatedFiles(t, fcmap) // Checking File Content

	// Testing File Deletion
	f2bDeleted := db.File{Title: "test_doc.copy.copy.txt", Password: "mypassword1"} // File to be deleted
	removeFile(t, f2bDeleted)                                                       // Deleting f2bDeleted ("test_doc.copy.copy.txt")

	// Check if deleted or not
	_, err := ioutil.ReadFile(path.Join(db.FileLoc, f2bDeleted.Title))
	if err == nil {
		t.Error("Error:", "Unable to Delete file")
	}
}

func copyFile(t *testing.T, file db.File, src string) (string, error) {
	np, err := file.CopyFile(src)
	if err != nil && err.Error() != "dup:err" {
		t.Error("Error:", err)
	}

	return np, nil
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

func checkCreatedFiles(t *testing.T, fcmap map[string]string) {
	for k, v := range fcmap {
		// filesl
		b, err := ioutil.ReadFile(path.Join(db.FileLoc, k))
		if err != nil {
			t.Error("Error:", err)
		}

		if string(b) != v {
			t.Error("Error:", errors.New("Content Mismatch, "+path.Join(db.FileLoc, k)))
		}
	}
}
