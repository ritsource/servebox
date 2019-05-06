package cmd

import (
	"fmt"

	db "github.com/ritwik310/servebox/database"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "r",
	Long:  `a`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var filetitle string
		var password string

		if len(args) == 0 {
			// if arguement not provided, ask in prompt
			fmt.Printf("File (File-path): ")
			fmt.Scanln(&filetitle) // Getting source path from the user
		} else {
			filetitle = args[0] // if arguement provided
		}

		// Getting Password
		fmt.Printf("Password: ")
		fmt.Scanln(&password) // Getting source path from the user

		f := db.File{Title: filetitle, Password: password}
		pw := db.Password{Title: filetitle, Password: password}

		// Removing file
		err := f.RemoveFile()
		if err != nil {
			fmt.Println("RemoveFileErr")
			return err
		}

		// Remove after deleting file
		err = pw.Remove()
		fmt.Println("RemoveErr")
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
