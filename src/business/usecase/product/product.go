package product

import (
	"context"
	productDom "go-clean/src/business/domain/product"
	"go-clean/src/business/entity"
)

type Interface interface {
	GetList(ctx context.Context, param entity.ProductParam) ([]entity.Product, error)
	Get(ctx context.Context, param entity.ProductParam) (entity.Product, error)
}

type product struct {
	product productDom.Interface
}

func Init(pd productDom.Interface) Interface {
	p := &product{
		product: pd,
	}

	return p
}

func (p *product) GetList(ctx context.Context, param entity.ProductParam) ([]entity.Product, error) {
	products, err := p.product.GetList(ctx, param)
	if err != nil {
		return products, err
	}

	return products, nil
}

func (p *product) Get(ctx context.Context, param entity.ProductParam) (entity.Product, error) {
	product, err := p.product.Get(ctx, param)
	if err != nil {
		return product, err
	}

	return product, nil
}
