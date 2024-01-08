package category

import (
	categoryDom "go-clean/src/business/domain/category"
	"go-clean/src/business/entity"
)

type Interface interface {
	GetList() ([]entity.Category, error)
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

func (c *category) GetList() ([]entity.Category, error) {
	categories, err := c.category.GetList()
	if err != nil {
		return categories, err
	}

	return categories, nil
}
