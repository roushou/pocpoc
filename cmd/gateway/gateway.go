package main

import (
	"context"
	"log"

	"github.com/roushou/pocpoc/internal/config"
	"github.com/roushou/pocpoc/internal/database"
	"github.com/roushou/pocpoc/internal/gateway"
	"github.com/roushou/pocpoc/internal/models"
	"github.com/roushou/pocpoc/internal/router"
	"gorm.io/gorm"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()
	db, err := database.NewDatabase(config.DatabaseName)
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}

	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}

	seedDatabase(db.Connection)

	router := router.NewRouter(db)

	gateway, err := gateway.NewGateway(router, gateway.WithAddr(config.GatewayAddr))
	if err != nil {
		log.Fatalf("failed to create gateway: %v", err)
	}

	gateway.Serve(ctx)
}

// seedDatabase seeds the database if no user records is found. For simplicity, it assumes other tables are empty or not based on that.
//
// Only for demo purposes
func seedDatabase(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Owner{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("Database already seeded. Skipping...")
		return nil
	}

	// Add owners
	owner1 := &models.Owner{Username: "Senku", PasswordHash: "strongpassword"}
	db.Create(owner1)

	// Add restaurants
	restaurant1 := &models.Restaurant{Name: "Sushi Den", OwnerID: owner1.ID}
	if err := db.Create(restaurant1).Error; err != nil {
		return err
	}

	// Add staff
	staff1 := &models.Staff{Username: "Taiju", PasswordHash: "1234", RestaurantID: restaurant1.ID}
	if err := db.Create(staff1).Error; err != nil {
		return err
	}

	// Add products
	product1 := &models.Product{
		Title:        "Wagyu Beef",
		Description:  "Tender wagyu beef",
		RestaurantID: restaurant1.ID,
	}
	if err := db.Create(product1).Error; err != nil {
		return err
	}

	return nil
}
