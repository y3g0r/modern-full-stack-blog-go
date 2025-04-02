package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	PORT string
}

func Load() Config {
	// Load the .env file
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading .env file")
	}

	// Set the default values
	// viper.SetDefault("PORT", "8080")

	// Read the values
	port := viper.GetString("PORT")
	return Config{
		PORT: port,
	}
}
