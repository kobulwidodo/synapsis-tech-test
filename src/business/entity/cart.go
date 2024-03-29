package entity

import "gorm.io/gorm"

const (
	StatusInCart = "in_cart"
	StatusUnpaid = "unpaid"
	StatusPaid   = "paid"
)

type Cart struct {
	gorm.Model
	UserID            uint
	ProductID         uint
	TransactionID     uint
	Qty               int
	Status            string
	FinalPricePerItem int
	TotalPriceNow     int64   `gorm:"-:all"`
	Product           Product `gorm:"-:all"`
}

type CartParam struct {
	ID            uint `uri:"cart_id"`
	UserID        uint
	ProductID     uint
	Status        string
	TransactionID uint
}

type CreateCartParam struct {
	ProductID uint `binding:"required"`
	Qty       int  `binding:"required"`
}

type UpdateCartParam struct {
	Qty               int
	Status            string
	TransactionID     uint
	FinalPricePerItem int
}
