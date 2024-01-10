package entity

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID      uint
	AddressShip string
	TotalPrice  int64
}

type CreateTransactionParam struct {
	AddressShip string `binding:"required"`
	PaymentID   int    `binding:"required"`
}

type TransactionParam struct {
	ID uint
}
