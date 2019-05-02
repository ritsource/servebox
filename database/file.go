package database

import (
	"io"
	"os"
	"path"
	"strings"
)

// File type represents each file in database
type File struct {
	Location string // Location can only be saved in password
	FileType string
	isPublic bool
	Password string
}

// IsExist checks if Filepath (fp) exist or not
// Returns path, if exist else returns "", returns "", err on error
func IsExist(fp string) (string, error) {
	if _, err := os.Stat(fp); err == nil {
		return fp, nil // File exists
	} else if os.IsNotExist(err) {
		return "", nil // File does not exists
	} else {
		return "", err
	}
}

// CopyData creates a destination file
// and copies data from a source file to the destination
func CopyData(src string, dst string) error {
	// Reading from Source file
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Creating new file
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copying Data from Source to Destination
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// GetFile serves the file if authenticated, returns http.FileSystem type
func (f File) GetFile() (string, error) {
	// Filepath, file's location
	fp := path.Join(FileLoc, f.Password, f.Location)
	return IsExist(fp)
}

// CopyFile copies a file from main source to servebox directory
func (f File) CopyFile(src string) (string, error) {
	// Filepath, file's location
	fp := path.Join(FileLoc, f.Password, f.Location)

	np, err := IsExist(fp)
	if err != nil {
		return "", err
	}

	// If path does not exist
	if np == "" {
		err = CopyData(src, fp)
		if err != nil {
			return "", err
		}
		return fp, nil
	}

	// If a file already a file exist
	// recursively try with new name by adding "-copy"

	nf := File(f)                          // First create a new File type
	nf.Location = genDupTitle(nf.Location) // Rename that &#!T

	return nf.CopyFile(src) // Recursively run the function
}

// Generates new filename that includes ".copy" before the extension,
// can be used to handle duplicate file
func genDupTitle(oldstr string) string {
	oldsl := strings.Split(oldstr, ".")
	newsl := append(oldsl[0:len(oldsl)-1], "copy", oldsl[len(oldsl)-1])
	return strings.Join(newsl, ".")
}
