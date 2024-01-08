package product

import (
	productDom "go-clean/src/business/domain/product"
	"go-clean/src/business/entity"
)

type Interface interface {
	GetList(param entity.ProductParam) ([]entity.Product, error)
	Get(param entity.ProductParam) (entity.Product, error)
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

func (p *product) GetList(param entity.ProductParam) ([]entity.Product, error) {
	products, err := p.product.GetList(param)
	if err != nil {
		return products, err
	}

	return products, nil
}

func (p *product) Get(param entity.ProductParam) (entity.Product, error) {
	product, err := p.product.Get(param)
	if err != nil {
		return product, err
	}

	return product, nil
}
