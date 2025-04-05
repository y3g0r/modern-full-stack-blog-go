package configs

import (
	"os"
)

type Config struct {
	PORT string
}

func Load() Config {
	return Config{
		PORT: getEnv("PORT", "3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
