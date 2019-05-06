package cmd

import (
	"fmt"

	db "github.com/ritwik310/servebox/database"
	"github.com/spf13/cobra"
)

// changePasswordCmd represents the changePassword command
var changePasswordCmd = &cobra.Command{
	Use:   "change-password",
	Short: "Changes the password of a file",
	RunE: func(cmd *cobra.Command, args []string) error {
		var newpassword string // Password input
		var filetitle string   // Filename (File's Title)

		if len(args) == 0 {
			// if arguement not provided, ask in prompt
			fmt.Printf("File: ")
			fmt.Scanln(&filetitle) // Getting File Title from the user
		} else {
			filetitle = args[0] // if arguement provided
		}

		fmt.Printf("New Password: ")
		fmt.Scanln(&newpassword) // Getting New Password path from the user

		// Reading password
		pw := db.Password{Title: filetitle}
		err := pw.Read() // Read passowrd just to check if it exist or not
		if err != nil {
			return err
		}

		// New Password struct
		newpw := db.Password{Title: pw.Title, Password: newpassword, FileName: pw.FileName}

		// Writing new password, it will replace the previous one
		err = newpw.Write()
		if err != nil {
			return err
		}

		// Printing teh results
		fmt.Print("\nSuccessfully Changed Password!\n\n")
		fmt.Println("From", pw.Password, "to", newpassword)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(changePasswordCmd)
}
