package usecase

import (
	"go-clean/src/business/domain"
	"go-clean/src/business/usecase/category"
	"go-clean/src/business/usecase/product"
	"go-clean/src/business/usecase/user"
	"go-clean/src/lib/auth"
)

type Usecase struct {
	User     user.Interface
	Category category.Interface
	Product  product.Interface
}

func Init(auth auth.Interface, d *domain.Domains) *Usecase {
	uc := &Usecase{
		User:     user.Init(d.User, auth),
		Category: category.Init(d.Category),
		Product:  product.Init(d.Product),
	}

	return uc
}
