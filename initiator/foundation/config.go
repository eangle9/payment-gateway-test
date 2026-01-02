package foundation

import (
	"log"

	"github.com/spf13/viper"
)

// InitConfig initializes viper to read configurations
func InitConfig() {
	// Set the file name of the configuration file
	viper.SetConfigFile(".env")

	// Automatically read system environment variables if the key doesn't exist in the file
	viper.AutomaticEnv()

	// Attempt to read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
}
