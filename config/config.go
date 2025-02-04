package config

import (
	"os"
)

type Config struct {
	Server struct {
		Port string
		Host string
	}

	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}

	ShortURL struct {
		Length  int
		BaseURL string
	}
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	cfg.Server.Port = getEnv("SERVER_PORT", "8080")
	cfg.Server.Host = getEnv("SERVER_HOST", "localhost")

	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "3306")
	cfg.Database.User = getEnv("DB_USER", "root")
	cfg.Database.Password = getEnv("DB_PASSWORD", "12345")
	cfg.Database.Name = getEnv("DB_NAME", "urlshortner")

	cfg.ShortURL.Length = 6
	cfg.ShortURL.BaseURL = "http://localhost:" + cfg.Server.Port

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
