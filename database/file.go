package database

import (
	"path"
)

// File type represents each file in database
type File struct {
	Location string // Location can only be saved in password
	Password string
}

// GetFile serves the file if authenticated, returns http.FileSystem type
func (f File) GetFile() string {
	// You need to go through authentication first
	return path.Join(FileLoc, f.Location)
}
