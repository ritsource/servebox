package database

import (
	"flag"
	"log"
	"os/user"
	"path"
)

// BaseLoc is the base directory for the application database
var BaseLoc string

// FileLoc is the directory for saving (copying) files
var FileLoc string

// PassLoc is the directory to save passwords
var PassLoc string

func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// Check if testing env or oot
	if flag.Lookup("test.v") == nil {
		BaseLoc = path.Join(usr.HomeDir, "servebox") // Default Base Directory
	} else {
		BaseLoc = path.Join("..", "test_servebox") // Base Directory in Testing Environment
	}

	FileLoc = path.Join(BaseLoc, "files")
	PassLoc = path.Join(BaseLoc, "passwords")
}
