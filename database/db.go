package database

import (
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

	BaseLoc = path.Join(usr.HomeDir, "servebox")
	FileLoc = path.Join(BaseLoc, "files")
	PassLoc = path.Join(BaseLoc, "passwords")
}
