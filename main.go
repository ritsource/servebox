package main

import (
	"fmt"
	"log"
	"os"
	"path"

	db "github.com/ritwik310/servebox/database"
	"github.com/ritwik310/servebox/server"
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
	f := db.File{
		Title:    "myfile.txt",
		Password: "mypassword0",
	}

	src := "/home/ritwik310/Downloads/test_doc.txt"

	np, err := testCopy(f, src)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("New Path:", np)

	npDir, npFile := path.Split(np)
	fmt.Println("npDir", npDir)
	fmt.Println("npFile", npFile)

	pw := db.Password{
		// Title:    "myfile.txt",
		Title:    npFile,
		Password: "mypassword0",
		FileName: npFile,
	}

	writePassword(pw)

	server.Start()
}

func testCopy(file db.File, src string) (string, error) {
	np, err := file.CopyFile(src)
	if err == nil {
		return np, err
	}

	if err.Error() == "dup:err" {
		// np, err = file.CopyFileRename(src, "newfile.txt")
		return file.CopyFileDup(src)
	}

	return "", err
}

func writePassword(pw db.Password) {
	err := pw.Write()
	if err != nil {
		log.Fatal(err)
	}
}

func readPassword() {
	pw := db.Password{
		Title:    "index3.txt",
		Password: "mypassword1",
		// FileName: "index3.txt",
	}

	err := pw.Read()
	if err == nil {
		fmt.Printf("%+v\n", pw)
		return
	}

	if err.Error() == "wrong:password" {
		fmt.Printf("Wrong Password!!!")
		return
	}

	log.Fatal(err)

}
