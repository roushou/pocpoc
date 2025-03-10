package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleOwner Role = "owner"
	RoleStaff Role = "staff"
)

type User struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username     string
	PasswordHash string
	Role         Role      `gorm:"not null"` // Roles: 'owner' or 'staff'
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	userID, err := uuid.NewV7()
	u.ID = userID
	return
}
