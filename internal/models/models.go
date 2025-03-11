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

type Owner struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username     string    `gorm:"not null; unique"`
	PasswordHash string    `gorm:"not null"`
}

func (o *Owner) BeforeCreate(tx *gorm.DB) (err error) {
	ID, err := uuid.NewV7()
	if err != nil {
		return
	}
	o.ID = ID
	return
}

type Staff struct {
	gorm.Model
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey"`
	RestaurantID uuid.UUID  `gorm:"type:uuid;not null"`
	Username     string     `gorm:"not null"`
	PasswordHash string     `gorm:"not null"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID;references:ID"`
}

func (s *Staff) BeforeCreate(tx *gorm.DB) (err error) {
	ID, err := uuid.NewV7()
	if err != nil {
		return
	}
	s.ID = ID
	return
}

type Restaurant struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	OwnerID uuid.UUID `gorm:"type:uuid;not null"`
	Name    string    `gorm:"not null"`
	Owner   Owner     `gorm:"foreignKey:OwnerID;references:ID"`
}

func (r *Restaurant) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	r.ID = id
	return
}
