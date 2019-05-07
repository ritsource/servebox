package cmd

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/fatih/color"
	"github.com/ritwik310/servebox/server"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a server on your local machine, that serves the files",
	RunE: func(cmd *cobra.Command, args []string) error {
		// fmt.Println("start called")
		var customport string
		var port string

		port = "6060"

		fmt.Printf("Custom Port [y/N]: ")
		fmt.Scanln(&customport)

		if strings.ToLower(customport) == "y" {
			fmt.Printf("Enter Port: ")
			fmt.Scanln(&port)
		}

		// Clear the terminal
		print("\033[H\033[2J")

		// WARNING Message
		color.Red("WARNING!, your connection is not private. To learn more, https://github.com/ritwik310/servebox#how-it-works")

		// Printing URL
		locIP := GetOutboundIP()
		color.Yellow("You can access the files at http://" + locIP.String() + ":" + port)

		return server.Start(port)
	},
}

// GetOutboundIP gets preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func init() {
	rootCmd.AddCommand(startCmd)
}
