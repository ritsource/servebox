package database

import (
	"errors"
	"io"
	"os"
)

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

// WriteFile does the writing for File methods, dries up the code
func WriteFile(fp string, src string) (string, error) {
	// Check if fp (filepath) exists or not
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

	// If a file already a file exist then return a "dup:err" error
	return "", errors.New("dup:err")
}
