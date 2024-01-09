package domain

import (
	"go-clean/src/business/domain/cart"
	"go-clean/src/business/domain/category"
	"go-clean/src/business/domain/product"
	"go-clean/src/business/domain/user"

	"gorm.io/gorm"
)

type Domains struct {
	User     user.Interface
	Category category.Interface
	Product  product.Interface
	Cart     cart.Interface
}

func Init(db *gorm.DB) *Domains {
	d := &Domains{
		User:     user.Init(db),
		Category: category.Init(db),
		Product:  product.Init(db),
		Cart:     cart.Init(db),
	}

	return d
}
