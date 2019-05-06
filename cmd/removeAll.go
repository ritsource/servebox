package cmd

import (
	"fmt"
	"os"

	db "github.com/ritwik310/servebox/database"
	"github.com/spf13/cobra"
)

// removeAllCmd represents the removeAll command
var removeAllCmd = &cobra.Command{
	Use:   "remove-all",
	Short: "Removes all files from the staging area",
	RunE: func(cmd *cobra.Command, args []string) error {
		// fmt.Println("removeAll called")
		err := os.RemoveAll(db.PassLoc)
		err = os.RemoveAll(db.FileLoc)
		if err == nil {
			fmt.Println("Successfully removed all Files and Passwords")
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(removeAllCmd)
}
