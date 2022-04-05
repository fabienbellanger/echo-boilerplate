package cli

import (
	server "github.com/fabienbellanger/echo-boilerplate"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "run",
	Short: "Start server",
	Long:  `Start server`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func startServer() {
	initConfig()

	// Start server
	// ------------
	server.Run()
}
