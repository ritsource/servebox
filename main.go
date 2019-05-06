package main

import (
	"log"
	"os"

	"github.com/ritwik310/servebox/cmd"
	db "github.com/ritwik310/servebox/database"
)

func init() {
	// Creates required directories, doesn't delete if already exist
	err := os.MkdirAll(db.FileLoc, os.ModePerm)
	err = os.MkdirAll(db.PassLoc, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	cmd.Execute()
}
