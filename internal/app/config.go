package app

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	DatabaseURL        string
	RedisURL           string
	AdminToken         string
	EncryptionKey      string
	ListenAddr         string
	RequestTimeout     time.Duration
	RateLimitPerMinute int
}

func LoadConfig() (Config, error) {
	c := Config{DatabaseURL: os.Getenv("DATABASE_URL"), RedisURL: os.Getenv("REDIS_URL"), AdminToken: os.Getenv("ADMIN_TOKEN"), EncryptionKey: os.Getenv("ENCRYPTION_KEY"), ListenAddr: env("LISTEN_ADDR", ":8080"), RequestTimeout: 90 * time.Second, RateLimitPerMinute: 60}
	if c.DatabaseURL == "" {
		return c, fmt.Errorf("DATABASE_URL is required")
	}
	if len(c.AdminToken) < 24 {
		return c, fmt.Errorf("ADMIN_TOKEN must contain at least 24 characters")
	}
	if len(c.EncryptionKey) < 24 {
		return c, fmt.Errorf("ENCRYPTION_KEY must contain at least 24 characters")
	}
	return c, nil
}
func env(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}
