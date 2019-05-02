package database

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

// Password type represents password file in the database
type Password struct {
	Title    string
	Password string
	FileName string
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

	// Writing file
	err := ioutil.WriteFile(path.Join(PassLoc, p.Title), b, 0777)
	if err != nil {
		return err
	}

	return nil
}

// Read reads password files, and populates the Password.FileName value
func (p *Password) Read() error {
	// Reading file
	b, err := ioutil.ReadFile(path.Join(PassLoc, p.Title))
	if err != nil {
		return err
	}

	// Extracting Password-File's data
	arr0 := bytes.Split(b, []byte("\n"))          // arr0 - array of each-line-data in the file
	bPwStr := arr0[0]                             // Password data (First line in the file)
	idxPwSp := bytes.IndexByte(bPwStr, byte(' ')) // Index of First " " - whitespace in First line (idxPwSp)

	// Check if p.Password is Not Correct
	if !bytes.Equal(bPwStr[idxPwSp+1:], []byte(p.Password)) {
		fmt.Println("Wrong Password!")
		return errors.New("wrong:password")
	}

	// If Password is Correct
	bFnStr := arr0[1]                             // Filename data (Second line in the file)
	idxFnSp := bytes.IndexByte(bFnStr, byte(' ')) // Index of First " " - whitespace in Second line (bFnStr)

	// Change FileName in p *Password
	p.FileName = string(bFnStr[idxFnSp+1:])

	return nil
}

// func CreatePassword
