package usecase

import (
	"go-clean/src/business/domain"
	"go-clean/src/business/usecase/cart"
	"go-clean/src/business/usecase/category"
	midtranstransaction "go-clean/src/business/usecase/midtrans_transaction"
	"go-clean/src/business/usecase/product"
	"go-clean/src/business/usecase/transaction"
	"go-clean/src/business/usecase/user"
	"go-clean/src/lib/auth"
)

type Usecase struct {
	User                user.Interface
	Category            category.Interface
	Product             product.Interface
	Cart                cart.Interface
	Transaction         transaction.Interface
	MidtransTransaction midtranstransaction.Interface
}

func Init(auth auth.Interface, d *domain.Domains) *Usecase {
	uc := &Usecase{
		User:                user.Init(d.User, auth),
		Category:            category.Init(d.Category),
		Product:             product.Init(d.Product),
		Cart:                cart.Init(d.Cart, auth, d.Product),
		Transaction:         transaction.Init(auth, d.Transaction, d.Cart, d.Product, d.Midtrans, d.MidtransTransaction),
		MidtransTransaction: midtranstransaction.Init(d.MidtransTransaction, d.Midtrans, d.Cart),
	}

	return uc
}
