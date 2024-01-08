package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	CategoryId  uint
	Name        string
	Description string
	Price       int
}

type ProductParam struct {
	ID         uint `uri:"product_id"`
	CategoryId uint `form:"category_id"`
}
