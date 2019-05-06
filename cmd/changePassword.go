package cmd

import (
	"fmt"

	"github.com/ritwik310/servebox/database"
	"github.com/spf13/cobra"
)

// changePasswordCmd represents the changePassword command
var changePasswordCmd = &cobra.Command{
	Use:   "change-password",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		pw := database.Password{Title: filetitle}
		err := pw.Read() // Read passowrd just to check if it exist or not
		if err != nil {
			return err
		}

		// New Password struct
		newpw := database.Password{Title: pw.Title, Password: newpassword, FileName: pw.FileName}

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// changePasswordCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// changePasswordCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
