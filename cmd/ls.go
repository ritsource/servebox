package cmd

import (
	"fmt"
	"os"
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
				pathSl := strings.Split(path, string(os.PathSeparator))[5:]
				relPath := string(os.PathSeparator) + strings.Join(pathSl, string(os.PathSeparator)) // relative path inside .servebox

				// Reading password
				pw := db.Password{Title: relPath}

				err := pw.Read()
				if err != nil {
					return err
				}

				fmt.Println("File: " + pw.Title + "\tPassword: " + pw.Password)
			}
			return nil
		})

		return err
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
