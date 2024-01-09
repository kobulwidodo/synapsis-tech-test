package cart

import (
	"context"
	cartDom "go-clean/src/business/domain/cart"
	productDom "go-clean/src/business/domain/product"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
)

type Interface interface {
	Create(ctx context.Context, cartInput entity.CreateCartParam) (entity.Cart, error)
}

type cart struct {
	cart    cartDom.Interface
	product productDom.Interface
	auth    auth.Interface
}

func Init(cd cartDom.Interface, auth auth.Interface, pd productDom.Interface) Interface {
	c := &cart{
		cart:    cd,
		auth:    auth,
		product: pd,
	}

	return c
}

func (c *cart) Create(ctx context.Context, cartInput entity.CreateCartParam) (entity.Cart, error) {
	result := entity.Cart{}

	user, err := c.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return result, err
	}

	product, err := c.product.Get(entity.ProductParam{
		ID: cartInput.ProductID,
	})
	if err != nil {
		return result, err
	}

	cartExist, _ := c.cart.Get(entity.CartParam{
		UserID:    user.User.ID,
		ProductID: product.ID,
		Status:    entity.StatusInCart,
	})

	if cartExist.ID != 0 {
		if err := c.cart.Update(entity.CartParam{
			UserID:    user.User.ID,
			ProductID: product.ID,
			Status:    entity.StatusInCart,
		}, entity.UpdateCartParam{
			Qty: cartExist.Qty + cartInput.Qty,
		}); err != nil {
			return cartExist, err
		}

		return cartExist, nil
	}

	result, err = c.cart.Create(entity.Cart{
		UserID:    user.User.ID,
		ProductID: product.ID,
		Qty:       cartInput.Qty,
		Status:    entity.StatusInCart,
	})
	if err != nil {
		return result, err
	}

	return result, nil
}
