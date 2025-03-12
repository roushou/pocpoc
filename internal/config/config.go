package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AllowedOrigins []string
	GatewayAddr    string
	DatabaseName   string
}

// LoadConfig loads all required configuration from environment variables.
func LoadConfig() (*Config, error) {
	allowedOriginsStr, ok := os.LookupEnv("ALLOWED_ORIGINS")
	if !ok {
		return nil, fmt.Errorf("environment variable 'DATABASE_NAME' not found")
	}
	allowedOrigins := strings.Split(allowedOriginsStr, ",")

	gatewayAddr, ok := os.LookupEnv("GATEWAY_ADDR")
	if !ok {
		return nil, fmt.Errorf("environment variable 'GATEWAY_ADDR' not found")
	}
	dbName, ok := os.LookupEnv("DATABASE_NAME")
	if !ok {
		return nil, fmt.Errorf("environment variable 'DATABASE_NAME' not found")
	}

	return &Config{
		AllowedOrigins: allowedOrigins,
		GatewayAddr:    gatewayAddr,
		DatabaseName:   dbName,
	}, nil
}
