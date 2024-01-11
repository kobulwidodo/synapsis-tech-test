package domain

import (
	"go-clean/src/business/domain/cart"
	"go-clean/src/business/domain/category"
	"go-clean/src/business/domain/midtrans"
	midtranstransaction "go-clean/src/business/domain/midtrans_transaction"
	"go-clean/src/business/domain/product"
	"go-clean/src/business/domain/transaction"
	"go-clean/src/business/domain/user"
	midtransSdk "go-clean/src/lib/midtrans"
	"go-clean/src/lib/redis"

	"gorm.io/gorm"
)

type Domains struct {
	User                user.Interface
	Category            category.Interface
	Product             product.Interface
	Cart                cart.Interface
	Midtrans            midtrans.Interface
	Transaction         transaction.Interface
	MidtransTransaction midtranstransaction.Interface
}

func Init(db *gorm.DB, m midtransSdk.Interface, redis redis.Interface) *Domains {
	d := &Domains{
		User:                user.Init(db),
		Category:            category.Init(db, redis),
		Product:             product.Init(db, redis),
		Cart:                cart.Init(db),
		Midtrans:            midtrans.Init(m),
		Transaction:         transaction.Init(db),
		MidtransTransaction: midtranstransaction.Init(db),
	}

	return d
}
