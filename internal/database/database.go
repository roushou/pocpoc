package database

import (
	"fmt"

	"github.com/roushou/pocpoc/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Connection *gorm.DB
}

func NewDatabase(databaseName string) (*Database, error) {
	connection, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return &Database{connection}, nil
}

func (db *Database) AutoMigrate() error {
	return db.Connection.AutoMigrate(
		&models.Owner{},
		&models.Restaurant{},
		&models.Staff{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
	)
}
