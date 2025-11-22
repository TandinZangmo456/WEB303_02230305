package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID     uint        `gorm:"not null"`
	Status     string      `gorm:"default:'pending'"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID      uint
	MenuItemID   uint
	MenuItemName string
	Quantity     uint
	Price        float64
}
