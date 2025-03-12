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

type Product struct {
	gorm.Model
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey"`
	RestaurantID uuid.UUID  `gorm:"type:uuid;not null"`
	Title        string     `gorm:"not null"`
	Description  string     `gorm:"not null"`
	UnitPrice    float64    `gorm:"not null"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID;references:ID"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	p.ID = id
	return
}

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusPrepared  OrderStatus = "prepared"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	gorm.Model
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey"`
	RestaurantID uuid.UUID   `gorm:"type:uuid;not null"`
	StaffID      uuid.UUID   `gorm:"type:uuid;not null"`
	TableNumber  string      `gorm:"not null"`
	Status       OrderStatus `gorm:"not null;default:pending"`
	TotalAmount  float64     `gorm:"not null;default:0.0"`
	Restaurant   Restaurant  `gorm:"foreignKey:RestaurantID;references:ID"`
	OrderItems   []OrderItem `gorm:"foreignKey:OrderID;references:ID"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	o.ID = id
	return
}

type OrderItem struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Quantity  uint32    `gorm:"not null"`
	Product   Product   `gorm:"foreignKey:ProductID;references:ID"`
}

func (o *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		return
	}
	o.ID = id
	return
}

func CalculateTotalAmount(items []OrderItem) float64 {
	totalAmount := 0.0
	for _, item := range items {
		totalAmount += float64(item.Quantity) * item.Product.UnitPrice
	}
	return totalAmount
}
