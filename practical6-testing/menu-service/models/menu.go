package models

import "gorm.io/gorm"

type MenuItem struct {
	gorm.Model
	Name        string  `gorm:"not null"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"not null"`
}
