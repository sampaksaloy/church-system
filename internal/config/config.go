package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	ServerPort string
	ServerHost string

	SessionSecret string

	AdminDefaultEmail    string
	AdminDefaultPassword string

	UploadMaxSizeMB int
	UploadDir       string
}

func Load() *Config {
	uploadMaxSize, _ := strconv.Atoi(getEnv("UPLOAD_MAX_SIZE_MB", "10"))

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "church_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		ServerPort: getEnv("SERVER_PORT", "8080"),
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),

		SessionSecret: getEnv("SESSION_SECRET", "church-secret-key-change-in-production"),

		AdminDefaultEmail:    getEnv("ADMIN_DEFAULT_EMAIL", "admin@church.com"),
		AdminDefaultPassword: getEnv("ADMIN_DEFAULT_PASSWORD", "Admin@123"),

		UploadMaxSizeMB: uploadMaxSize,
		UploadDir:       getEnv("UPLOAD_DIR", "static/images/uploads"),
	}
}

func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
