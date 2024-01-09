package entity

import "gorm.io/gorm"

const (
	StatusInCart = "in_cart"
	StatusUnpaid = "unpaid"
	StuatusPaid  = "paid"
)

type Cart struct {
	gorm.Model
	UserID            uint
	ProductID         uint
	TransactionID     uint
	Qty               int
	Status            string
	FinalPricePerItem uint
}

type CartParam struct {
	ID        uint
	UserID    uint
	ProductID uint
	Status    string
}

type CreateCartParam struct {
	ProductID uint `binding:"required"`
	Qty       int  `binding:"required"`
}

type UpdateCartParam struct {
	Qty int
}
