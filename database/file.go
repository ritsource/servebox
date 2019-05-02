package database

import (
	"os"
	"path"
)

// File type represents each file in database
type File struct {
	Location string // Location can only be saved in password
	FileType string
	isPublic bool
	Password string
}

// GetFile serves the file if authenticated, returns http.FileSystem type
func (f File) GetFile() (string, error) {
	// Filepath, file's location
	fp := path.Join(FileLoc, f.Password, f.Location)

	// Check if Filepath (fp) exist or not
	if _, err := os.Stat(fp); err == nil {
		return fp, nil // File exists
	} else if os.IsNotExist(err) {
		return "", nil // File does not exists
	} else {
		return "", err
	}

	// return path.Join(FileLoc, f.Password, f.Location)
}
