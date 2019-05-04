package database

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Password type represents password file in the database
type Password struct {
	Title    string
	Password string
	FileName string
}

// GetFileName reads password files, and populates the Password.FileName value
func (p *Password) GetFileName() error {
	// New variable to store Query Secrets
	p2 := Password{Title: p.Title}
	err := p2.Read()
	if err != nil {
		return err
	}

	// If Query secret (Password) matches given Password
	if p2.Password != p.Password {
		return errors.New("wrong:password")
	}

	// Change FileName in p *Password
	p.FileName = p2.FileName

	return nil
}

// Write writes a file that contains Password.Password and Password.FileName
func (p Password) Write() error {
	// Checking if password is valid or not, check if contains "\n"
	if strings.IndexByte(p.Password, byte('\n')) != -1 {
		return errors.New("Password is not acceptable, contains \"\\n\"")
	}

	// All three fields are required, (Unnecrssary)
	if p.Title == "" || p.Password == "" || p.FileName == "" {
		return errors.New("All three fields are required, Title, Password and Filename")
	}

	// Formatting Password-File Data
	bPwSl := bytes.Join([][]byte{[]byte("Password"), []byte(p.Password)}, []byte(" ")) // Password Data
	bFnSl := bytes.Join([][]byte{[]byte("FileName"), []byte(p.FileName)}, []byte(" ")) // Filename Dara
	b := bytes.Join([][]byte{bPwSl, bFnSl}, []byte("\n"))                              // The Whole PasswordFile Dara

	// Creating directory for Password-File
	fpDir, _ := path.Split(p.Title)                            // Directory name of p.Title
	err := os.MkdirAll(path.Join(PassLoc, fpDir), os.ModePerm) // Creating directory
	if err != nil {
		return err
	}

	// Writing file
	err = ioutil.WriteFile(path.Join(PassLoc, p.Title), b, 0777)
	if err != nil {
		return err
	}

	return nil
}

// Read populates all the values of password given it's Title
func (p *Password) Read() error {
	// Check if Title Provided
	if p.Title == "" {
		return errors.New("p.Title field required in Password")
	}

	// Reading file
	b, err := ioutil.ReadFile(path.Join(PassLoc, p.Title))
	if err != nil {
		return err
	}

	// Extracting Password-File's data
	arr := bytes.Split(b, []byte("\n")) // arr - array of each-line-data in the file

	bPwStr := arr[0]                              // Password data (First line in the file)
	idxPwSp := bytes.IndexByte(bPwStr, byte(' ')) // Index of First " " - whitespace in First line (idxPwSp)
	p.Password = string(bPwStr[idxPwSp+1:])       // Setting pw.Password

	bFnStr := arr[1]                              // Filename data (Second line in the file)
	idxFnSp := bytes.IndexByte(bFnStr, byte(' ')) // Index of First " " - whitespace in Second line (bFnStr)
	p.FileName = string(bFnStr[idxFnSp+1:])       // Setting pw.Filename

	return err
}

// Remove removes a password file
func (p Password) Remove() error {
	// Check if Related exist
	fp, err := IsExist(path.Join(FileLoc, p.FileName))
	if err != nil {
		return err
	}

	// if corresponding is already deleted or not
	if fp != "" {
		return errors.New("password cannot be deleted, cause corresponding file exists")
	}

	// Removing file
	err = os.Remove(path.Join(PassLoc, p.Title))
	return err
}
