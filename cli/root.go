package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "Echo Boilerplate",
	Short:   "A Echo boilerplate with GORM",
	Long:    "A Echo boilerplate with GORM",
	Version: version,
}

// Execute starts CLI
func Execute() error {
	return rootCmd.Execute()
}

// initConfig initializes configuration from config file.
func initConfig() error {
	viper.SetConfigFile(".env")
	return viper.ReadInConfig()
}
