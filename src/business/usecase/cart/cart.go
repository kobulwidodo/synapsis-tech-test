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
	GetList(ctx context.Context) ([]entity.Cart, error)
	Delete(ctx context.Context, param entity.CartParam) error
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

func (c *cart) GetList(ctx context.Context) ([]entity.Cart, error) {
	result := []entity.Cart{}

	user, err := c.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return result, err
	}

	result, err = c.cart.GetList(entity.CartParam{
		UserID: user.User.ID,
		Status: entity.StatusInCart,
	})
	if err != nil {
		return result, err
	}

	mapProductIDs := make(map[uint]bool)
	for _, c := range result {
		mapProductIDs[c.ProductID] = true
	}

	productIDs := []uint{}
	for id := range mapProductIDs {
		productIDs = append(productIDs, id)
	}

	products, err := c.product.GetListByID(productIDs)
	if err != nil {
		return result, err
	}

	productsMap := make(map[uint]entity.Product)
	for _, p := range products {
		productsMap[p.ID] = p
	}

	for i, c := range result {
		result[i].TotalPriceNow = int64(c.Qty * productsMap[c.ProductID].Price)
		result[i].Product = productsMap[c.ProductID]
	}

	return result, nil
}

func (c *cart) Delete(ctx context.Context, param entity.CartParam) error {
	user, err := c.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return err
	}

	if err := c.cart.Delete(entity.CartParam{
		ID:     param.ID,
		UserID: user.User.ID,
		Status: entity.StatusInCart,
	}); err != nil {
		return err
	}

	return nil
}
