package database_test

import (
	"path"
	"path/filepath"
	"testing"

	db "github.com/ritwik310/servebox/database"
)

func TestPasswordWriteReadAndRemove(t *testing.T) {
	filename := "test_pass.txt"                  // Name of file (Title)
	password := "mypassword1"                    // Password for File
	src := path.Join(SourceDir, "test_pass.txt") // Source File Path

	f := db.File{Title: filename} // File Struct, containing Title and password

	// Copying File from source
	np, err := CopyFile(t, f, src) // copyFile function takes care of error
	if err != nil {
		t.Error(err)
	}

	npFile := filepath.Base(np) // Geting directory-name and file-name from np (new filepath)
	pw := db.Password{
		Title:    npFile,
		Password: password,
		FileName: npFile,
	} // New Password Struct

	// Writing File
	writePassword(t, pw) // Password File Write Test

	// Reading File
	pwR := db.Password(pw) // New Password sStruct for Read test
	pwR.FileName = ""      // Set filename to Zero value, will check if query populates it

	readPassword(t, &pwR) // Password File Read Test
	if pwR.FileName != pw.FileName {
		t.Error("Password ouery doesn't match, original \"FileName\"")
	}

	// Reading with Wrong Password
	pwW := db.Password(pw)                    // Password Struct for Wrong Password Check
	pwW.Password = pwW.Password + "blah_blah" // Rewrite password

	readPasswordWrong(t, pwW) // Password File Wrong Password Read Test

	removePassword(t, f, pw) // Tests pw.Remove method
}

func writePassword(t *testing.T, pw db.Password) {
	err := pw.Write()
	if err != nil {
		t.Error(err)
	}
}

func readPassword(t *testing.T, pw *db.Password) {
	err := pw.GetFileName()
	if err == nil {
		return
	}

	t.Error(err)
}

func readPasswordWrong(t *testing.T, pw db.Password) {
	err := pw.GetFileName()
	if err == nil {
		t.Error("Not returning error on wrong password")
		return
	}

	if err.Error() == "wrong:password" {
		return
	}

	t.Error("Not returning \"wrong:password\" error, on wrong password")
}

func removePassword(t *testing.T, f db.File, pw db.Password) {
	err := pw.Remove() // Remove without deleting file
	if err == nil {
		t.Error("Should have thrown Error, if file not deleted beforehand")
	}

	err = f.RemoveFile() // Removing file
	if err != nil {
		t.Error(err)
	}

	err = pw.Remove() // Remove after deleting file
	if err != nil {
		t.Error(err)
	}
}
