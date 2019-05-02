package database

import (
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

// Utils here ...
// IsExist, CopyData, WriteFile, HandlePassDir

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

// HandlePassDir creates File.Password directory if does not exist
func HandlePassDir(pwPath string) error {
	// pwp, err := IsExist(pwPath) // Password path (pwp), check if exist
	err := os.MkdirAll(pwPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// File type represents each file in database
type File struct {
	Title    string // Location can only be saved in password
	FileType string
	isPublic bool
	Password string
}

// File Methods ...
// GetFile, CopyFile, CopyFileDup, CopyFileRename, RemoveFile

// GetFile serves the file if authenticated, returns http.FileSystem type
func (f File) GetFile() (string, error) {
	// Filepath, file's location
	fp := path.Join(FileLoc, f.Password, f.Title)

	// HandlePassDir takes care of password directory, creates if does not exist
	HandlePassDir(path.Join(FileLoc, f.Password))
	return IsExist(fp)
}

// CopyFile copies a file from main source to servebox directory
func (f File) CopyFile(src string) (string, error) {
	// Filepath, file's location
	fp := path.Join(FileLoc, f.Password, f.Title)

	HandlePassDir(path.Join(FileLoc, f.Password)) // To handle PasswordDir
	return WriteFile(fp, src)                     // WriteFile writes teh file, and handle errors (To dry up the code)
}

// CopyFileDup copies file data by recursively
// adding ".copy" to the title if filepath exist
func (f File) CopyFileDup(src string) (string, error) {
	nf := File(f)                    // First create a new File type
	nf.Title = genDupTitle(nf.Title) // Rename that Title &#!T

	fp := path.Join(FileLoc, nf.Password, nf.Title) // New Path that includes ".copy"

	HandlePassDir(path.Join(FileLoc, f.Password)) // To handle PasswordDir

	np, err := WriteFile(fp, src) // Writes file
	if err == nil {
		return np, nil // If no error then return expected returns
	}

	// If error refers that, this filename exists too, try again & again & again
	if err.Error() == "dup:err" {
		return nf.CopyFileDup(src) // Recursively run the function
	}

	// If some other error
	return "", err
}

// Generates new filename that includes ".copy" before the extension,
// can be used to handle duplicate file
func genDupTitle(oldstr string) string {
	oldsl := strings.Split(oldstr, ".")
	newsl := append(oldsl[0:len(oldsl)-1], "copy", oldsl[len(oldsl)-1])
	return strings.Join(newsl, ".")
}

// CopyFileRename copies teh file data with new name
func (f File) CopyFileRename(src string, nTit string) (string, error) {
	nf := File(f)   // First create a new File type
	nf.Title = nTit // Rename that Title &#!T, this time also

	fp := path.Join(FileLoc, nf.Password, nf.Title) // New Path

	np, err := WriteFile(fp, src) // Writes file
	if err == nil {
		return np, nil // If no error then return expected returns
	}

	// If error refers that, a file of that name already a file exist then return a "dup:err" error
	if err.Error() == "dup:err" {
		// return "", errors.New("dup:err")
		return "", err
	}

	return "", err // If some other error

	// NOTE: I could have just returned the error, but this is less confiusing and easy to understand
}

// RemoveFile deletes a file
func (f File) RemoveFile() error {
	// filepath
	fp := path.Join(FileLoc, f.Password, f.Title)

	// Removing file
	err := os.Remove(fp)
	if err != nil {
		return err
	}

	return nil
}
