package database_test

import (
	"testing"

	db "github.com/ritwik310/servebox/database"
)

func TestWriteAndRead(t *testing.T) {
	pw := db.Password{
		Title:    "test_pass.txt",
		Password: "mypassword1",
		FileName: "test_pass.txt",
	}

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
}

func writePassword(t *testing.T, pw db.Password) {
	err := pw.Write()
	if err != nil {
		t.Error(err)
	}
}

func readPassword(t *testing.T, pw *db.Password) {
	err := pw.Read()
	if err == nil {
		return
	}

	t.Error(err)
}

func readPasswordWrong(t *testing.T, pw db.Password) {
	err := pw.Read()
	if err == nil {
		t.Error("Not returning error on wrong password")
		return
	}

	if err.Error() == "wrong:password" {
		return
	}

	t.Error("Not returning \"wrong:password\" error, on wrong password")
}
