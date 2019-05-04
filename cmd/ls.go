package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	db "github.com/ritwik310/servebox/database"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "ls",
	Long:  "ls",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Reading db.FileLoc tree
		err := filepath.Walk(db.FileLoc, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Check if it's a file
			if !info.IsDir() {
				var printMsg string

				pathSl := strings.Split(path, string(os.PathSeparator))[5:]
				relPath := string(os.PathSeparator) + strings.Join(pathSl, string(os.PathSeparator)) // relative path inside .servebox
				printMsg = "File: " + relPath

				// Reading password
				pw, err := readPassword(relPath)
				if err != nil {
					return err
				}

				printMsg = printMsg + "\tPassword: " + pw

				fmt.Println(printMsg)
			}
			return nil
		})

		return err
	},
}

// readPassword Password from Password-file
func readPassword(relPath string) (string, error) {
	// Reading file
	b, err := ioutil.ReadFile(path.Join(db.PassLoc, relPath))
	if err != nil {
		return "", err
	}

	// Extracting Password-File's data
	arr0 := bytes.Split(b, []byte("\n"))          // arr0 - array of each-line-data in the file
	bPwStr := arr0[0]                             // Password data (First line in the file)
	idxPwSp := bytes.IndexByte(bPwStr, byte(' ')) // Index of First " " - whitespace in First line (idxPwSp)

	return string(bPwStr[idxPwSp+1:]), nil
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
