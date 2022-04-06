package cli

import (
	"time"

	"github.com/fabienbellanger/echo-boilerplate/db"
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

func initConfigLoggerDatabase(initDB bool) (database *db.DB, err error) {
	// Configuration initialization
	// ----------------------------
	if err = initConfig(); err != nil {
		return nil, err
	}

	// Database connection
	// -------------------
	if initDB {
		database, err = db.New(&db.DatabaseConfig{
			Driver:          viper.GetString("DB_DRIVER"),
			Host:            viper.GetString("DB_HOST"),
			Username:        viper.GetString("DB_USERNAME"),
			Password:        viper.GetString("DB_PASSWORD"),
			Port:            viper.GetInt("DB_PORT"),
			Database:        viper.GetString("DB_DATABASE"),
			Charset:         viper.GetString("DB_CHARSET"),
			Collation:       viper.GetString("DB_COLLATION"),
			Location:        viper.GetString("DB_LOCATION"),
			MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
			MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
			ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME") * time.Hour,
		})
		if err != nil {
			return nil, err
		}
	}

	return database, err
}
