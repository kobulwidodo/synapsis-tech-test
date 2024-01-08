package category

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	GetList() ([]entity.Category, error)
}

type cateogry struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	c := &cateogry{
		db: db,
	}

	return c
}

func (c *cateogry) GetList() ([]entity.Category, error) {
	categories := []entity.Category{}
	if err := c.db.Find(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}
