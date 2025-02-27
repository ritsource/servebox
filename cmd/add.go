package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	db "github.com/ritwik310/servebox/database"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a file to the staging area, ready to be shared publically",
	Long: `You can provide file source in the commarnd as the first arguement
Example, "servebox add ../test/example.txt"

Add command ask you to provide a "password", which is ironically meaningless because
the current version of servebox doesn't support TLS. More about that here
https://github.com/ritwik310/servebox/blob/master/README.md
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var replace bool     // File to be replaced or not (If already there)
		var arg0 string      // First arguement, path to the source file
		var src string       // Path to the source file
		var password string  // Password input
		var filetitle string // Filename (File's Title)
		var newpath string   // Path of the copied file (np)
		var err error

		// Reading Flags
		replace, err = cmd.Flags().GetBool("replace") // Reading "replace" flag
		if err != nil {
			return err
		}

		if len(args) == 0 {
			// if arguement not provided, ask in prompt
			fmt.Printf("Source Path: ")
			fmt.Scanln(&arg0) // Getting source path from the user
		} else {
			arg0 = args[0] // if arguement provided
		}

		// Checks for a valid file source
		arg0, err = checkSrc(arg0)
		if err != nil {
			fmt.Print("\n")
			return err
		}

		// Setting src to Absolete path of "arg0"
		src, err = filepath.Abs(arg0)
		if err != nil {
			fmt.Print("\n")
			return err
		}

		// Getting password from the user
		password, err = handlePasswordIn("Enter Password: ") // handlePasswordIn handles password input & validation
		if err != nil {
			return err
		}

		// Copying, Renaming & Replacing File
		filetitle = src                                 // Setting filetitle to src (Cause have to map out the whole fs from "/" (absolute))
		file := db.File{Title: src, Password: password} // New File struct

		newpath, err = file.CopyFile(src) // Copying source file

		if err != nil && err.Error() == "dup:err" {
			// If filetitle already exist
			if replace {
				// if user requested for existing file replacement (-r)
				err = handleReplace(src, &newpath, file) // handleReplace handles file replacing
				if err != nil {
					return err
				}
			} else {
				// If user doesn't wanna replace the previous one, create duplicate
				err = handleDuplicate(true, src, &filetitle, &newpath, &file) // handleDuplicate handles renaming
				if err != nil {
					return err
				}
			}
		} else if err != nil {
			// If other error
			return err
		}

		// Saving Password
		pw := db.Password{Title: filetitle, Password: password, FileName: filetitle} // New Password struct

		// Writing Password
		err = pw.Write()
		if err != nil {
			return err
		}

		// Printing teh results
		fmt.Print("\nSuccessfully Added!\n\n")
		fmt.Println("Source Location:", src)
		fmt.Println("Copy Location:", newpath)

		return nil
	},
}

// handlePasswordIn handles password input and validation (!= "")
func handlePasswordIn(label string) (string, error) {
	var passw string

	fmt.Print(label)
	fmt.Scanln(&passw) // Getting input

	if passw == "" {
		// If input == "", ask again
		return handlePasswordIn("Enter Password (Non-Optional): ")
	}

	return passw, nil
}

// handleReplace deletes the old file, and replaces it with the new file data
func handleReplace(src string, newpath *string, file db.File) error {
	// NOTE: handleReplace will remove file, but don't have to be worrid about password file deletion
	// cause while writing the password file later in the CMD-function we are gonnarewrite it anyway

	err := file.RemoveFile() // Remove the existing file
	if err != nil {
		return err
	}

	*newpath, err = file.CopyFile(src) // Copy the new Data
	if err != nil {
		return err
	}

	return nil
}

// handleDuplicate takes care of file renaming on duplicate copy
func handleDuplicate(askforauto bool, src string, filetitleP *string, newpathP *string, fileP *db.File) error {
	fmt.Printf("++ CONFLICT ++: filename \"" + *filetitleP + "\" already exist.\n")

	var autoRn string // Want-Auto-Rename input
	var err error
	fnDir, _ := path.Split(src) // fnDir is the directory that source is inside

	// Only want to ask for auto rename, indeally for the first time only
	if askforauto {
		fmt.Print("Rename or Auto-Rename (r/A): ")
		fmt.Scanln(&autoRn)
	}

	if askforauto && strings.ToLower(autoRn) != "r" {
		// When user askes for auto renaming
		*newpathP, err = fileP.CopyFileDup(src)
		if err != nil {
			return err
		}

		*filetitleP = path.Join(fnDir, filepath.Base(*newpathP))

	} else {
		var fnIn string // filetitle input

		fmt.Print("Enter Filename: ")
		fmt.Scanln(&fnIn) // Getting fnIn

		*filetitleP = path.Join(fnDir, fnIn)

		*newpathP, err = fileP.CopyFileRename(src, *filetitleP)

		if err != nil && err.Error() == "dup:err" {
			// If manual name is also duplicate, run again (this time dont askforauto)
			return handleDuplicate(false, src, filetitleP, newpathP, fileP)
		} else if err != nil {
			fmt.Printf("\n")
			return err
		}

	}

	return nil
}

// checkSrc checks if a file is valid or not, given its path
func checkSrc(input string) (string, error) {
	info, err := os.Stat(input)

	if err == nil {
		// If file path exist
		if info.IsDir() {
			return "", errors.New("Path " + input + " is a directory, not file") // If filepath is a directory
		}
		return input, nil // If everything's fine

	} else if os.IsNotExist(err) {
		return "", errors.New("File path " + input + " doesn't exist") // If file path does not exist

	} else {
		return "", errors.New("File path " + input + " is invalid") // if unknown error

	}

}

func init() {
	addCmd.Flags().BoolP("replace", "r", false, "Replace with the new copy if file is already there") // FLags

	rootCmd.AddCommand(addCmd)
}
