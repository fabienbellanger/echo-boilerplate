package cli

import (
	"log"

	server "github.com/fabienbellanger/echo-boilerplate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "run",
	Short: "Start Web Server",
	Long:  `Start Web Server`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func startServer() {
	// Configuration initialization
	// ----------------------------
	logger, db, err := initConfigLoggerDatabase(true, true)
	if err != nil {
		log.Fatalln(err)
	}

	// Database migrations
	// -------------------
	if viper.GetBool("GORM_AUTOMIGRATIONS") {
		db.MakeMigrations()
	}

	// Start server
	// ------------
	server.Run(logger, db)
}
