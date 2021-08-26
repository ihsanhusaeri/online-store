package entity

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserId            uint        `json:"userId"`
	UserName          string      `json:"userName"`
	TotalPay          float64     `json:"totalPay"`
	Status            string      `json:"status" gorm:"default:cart"` //car sebagai default value untuk status
	CheckoutExpiredAt *time.Time  `json:"checkoutExpiredAt"`
	OrderItems        []OrderItem `json:"orderItems"`
}

type OrderItem struct {
	gorm.Model
	ItemId    uint    `json:"itemId"`
	ItemName  string  `json:"itemName"`
	ItemPrice float64 `json:"itemPrice"`
	ItemQty   uint    `json:"itemQty"`
	OrderId   uint    `json:"orderId"`
}
