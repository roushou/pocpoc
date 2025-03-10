package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username     string
	PasswordHash string
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	userID, err := uuid.NewV7()
	u.ID = userID
	return
}
