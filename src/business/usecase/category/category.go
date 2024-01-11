package category

import (
	"context"
	categoryDom "go-clean/src/business/domain/category"
	"go-clean/src/business/entity"
)

type Interface interface {
	GetList(ctx context.Context) ([]entity.Category, error)
}

type category struct {
	category categoryDom.Interface
}

func Init(cd categoryDom.Interface) Interface {
	c := &category{
		category: cd,
	}

	return c
}

func (c *category) GetList(ctx context.Context) ([]entity.Category, error) {
	categories, err := c.category.GetList(ctx)
	if err != nil {
		return categories, err
	}

	return categories, nil
}
