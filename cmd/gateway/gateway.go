package main

import (
	"context"
	"log"

	"github.com/roushou/pocpoc/internal/config"
	"github.com/roushou/pocpoc/internal/gateway"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()
	db, err := gorm.Open(sqlite.Open(config.DatabaseName), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	gateway, err := gateway.NewGateway(db, gateway.WithAddr(config.GatewayAddr))
	if err != nil {
		log.Fatalf("failed to create gateway: %v", err)
	}

	gateway.Serve(ctx)
}
