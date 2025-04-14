package configs

import (
	"os"
)

type Config struct {
	PORT         string
	CLERK_SK     string
	DATABASE_URL string
}

func Load() Config {
	return Config{
		PORT:         getEnv("PORT", "3000"),
		CLERK_SK:     getRequiredEnv("CLERK_SECRET_KEY"),
		DATABASE_URL: getEnv("DATABASE_URL", "DATABASE_URL"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getRequiredEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic("Required configuration key is missing: " + key)
}
