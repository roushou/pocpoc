package config

import (
	"fmt"
	"os"
)

type Config struct {
	GatewayAddr  string
	DatabaseName string
}

// LoadConfig loads all required configuration from environment variables.
func LoadConfig() (*Config, error) {
	gatewayAddr, ok := os.LookupEnv("GATEWAY_ADDR")
	if !ok {
		return nil, fmt.Errorf("environment variable 'GATEWAY_ADDR' not found")
	}
	dbName, ok := os.LookupEnv("DATABASE_NAME")
	if !ok {
		return nil, fmt.Errorf("environment variable 'DATABASE_NAME' not found")
	}

	return &Config{
		GatewayAddr:  gatewayAddr,
		DatabaseName: dbName,
	}, nil
}
