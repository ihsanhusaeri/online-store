package entity

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock uint    `gorm:"check:stock >= 0" json:"stock"`
}
