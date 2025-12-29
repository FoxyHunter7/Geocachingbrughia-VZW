package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Port         string
	Env          string
	DatabasePath string
	DataDir      string
	FrontendURL  string
	JWT          JWTConfig
	SMTP         SMTPConfig
	ReminderDays int
	CORSOrigins  []string
}

type JWTConfig struct {
	Secret      string
	ExpiryHours int
}

type SMTPConfig struct {
	Host              string
	Port              int
	User              string
	Pass              string
	From              string
	NotificationEmail string
}

func Load() *Config {
	dbPath := getEnv("DATABASE_PATH", "./data/geocaching.db")
	dataDir := filepath.Dir(dbPath)

	return &Config{
		Port:         getEnv("PORT", "8080"),
		Env:          getEnv("ENV", "development"),
		DatabasePath: dbPath,
		DataDir:      dataDir,
		FrontendURL:  getEnv("FRONTEND_URL", "http://localhost:5173"),
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "change-me-in-production"),
			ExpiryHours: getEnvInt("JWT_EXPIRY_HOURS", 24),
		},
		SMTP: SMTPConfig{
			Host:              getEnv("SMTP_HOST", ""),
			Port:              getEnvInt("SMTP_PORT", 587),
			User:              getEnv("SMTP_USER", ""),
			Pass:              getEnv("SMTP_PASS", ""),
			From:              getEnv("SMTP_FROM", ""),
			NotificationEmail: getEnv("NOTIFICATION_EMAIL", ""),
		},
		ReminderDays: getEnvInt("REMINDER_DAYS", 3),
		CORSOrigins:  strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173"), ","),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// Validate checks for security issues in production configuration
func (c *Config) Validate() error {
	if c.IsProduction() {
		if c.JWT.Secret == "change-me-in-production" {
			return errors.New("JWT_SECRET must be set to a secure value in production")
		}
		if len(c.JWT.Secret) < 32 {
			return errors.New("JWT_SECRET must be at least 32 characters in production")
		}
	}
	return nil
}
