package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	db "github.com/ritwik310/servebox/database"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A",
	Long:  `A`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var src string
		var password string
		var filename string
		var newpath string // (np)
		var err error

		if len(args) > 0 {
			src, err = checkSrc(args[0])
			if err != nil {
				// fmt.Printf("\n")
				return err
			}
		} else {
			srcIn, err := askInput("Source Path: ") // Source Input
			if err != nil {
				fmt.Printf("\n")
				return err
			}

			src, err = checkSrc(srcIn[:len(srcIn)-1]) // NOTE: srcIn includes "\n" in the end
			if err != nil {
				fmt.Printf("\n")
				return err
			}
		}

		password, err = handlePasswordIn("Enter Password: ")
		if err != nil {
			return err
		}

		_, srcFn := path.Split(src) // Reading the filename from src
		filename = srcFn            // Setting filename to srcFn

		file := db.File{Title: filename, Password: password} // New File Struct

		// fmt.Printf("%+v\n", f)

		newpath, err = file.CopyFile(src)

		if err != nil && err.Error() == "dup:err" { // TODO: Replace the file option

			newpath, err = handleDuplicate(true, src, &filename, &newpath, &file)
			if err != nil {
				return err
			}

		} else if err != nil {
			return err
		}

		_, npFile := path.Split(newpath)

		pw := db.Password{Title: npFile, Password: password, FileName: npFile}

		// fmt.Printf("%+v\n", pw)

		err = pw.Write()
		if err != nil {
			return err
		}

		fmt.Println("newpath:", newpath)
		fmt.Print("\n")
		fmt.Println("Source Location:\t", newpath)
		fmt.Println("Copy Location:\t", newpath)

		fmt.Println("add called")

		return nil
	},
}

func handlePasswordIn(label string) (string, error) {
	pwIn, err := askInput(label) // Password for File
	if err != nil {
		fmt.Printf("\n")
		return "", err
	}

	passw := pwIn[:len(pwIn)-1]

	if passw == "" {
		return handlePasswordIn("Enter Password (Non-Optional): ")
	}

	return passw, nil
}

func handleDuplicate(askforauto bool, src string, filename *string, newpath *string, file *db.File) (string, error) {
	fmt.Printf("++ CONFLICT ++: filename \"" + *filename + "\" already exist.\n")

	var autoRn string
	var err error

	if askforauto {
		rnIn, err := askInput("Rename or Auto-Rename (r/A): ")
		if err != nil {
			fmt.Printf("\n")
			return "", err
		}

		autoRn = rnIn[:len(rnIn)-1]
	}

	if askforauto && strings.ToLower(autoRn) != "r" {
		*newpath, err = file.CopyFileDup(src)
		if err != nil {
			return "", err
		}
	} else {

		fnIn, err := askInput("Enter Filename: ") // Password for File
		if err != nil {
			fmt.Printf("\n")
			return "", err
		}

		*filename = fnIn[:len(fnIn)-1]

		*newpath, err = file.CopyFileRename(src, *filename)

		if err != nil && err.Error() == "dup:err" {
			// askforauto = false
			return handleDuplicate(false, src, filename, newpath, file)
		} else if err != nil {
			fmt.Printf("\n")
			return "", err
		}

	}

	return *newpath, nil
}

func askInput(label string) (string, error) {
	reader := bufio.NewReader(os.Stdin) // New Reader
	fmt.Print(label)
	return reader.ReadString('\n') // Reading string until '\n'
}

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
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
