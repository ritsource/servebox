package cmd

import (
	"github.com/ritwik310/servebox/server"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the server",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("start called")
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
