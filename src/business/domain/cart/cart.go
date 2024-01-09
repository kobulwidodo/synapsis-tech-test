package cart

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(cart entity.Cart) (entity.Cart, error)
	GetList(param entity.CartParam) ([]entity.Cart, error)
	Get(param entity.CartParam) (entity.Cart, error)
	Update(selectParam entity.CartParam, updateParam entity.UpdateCartParam) error
}

type cart struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	c := &cart{
		db: db,
	}

	return c
}

func (c *cart) Create(cart entity.Cart) (entity.Cart, error) {
	if err := c.db.Create(&cart).Error; err != nil {
		return cart, err
	}

	return cart, nil
}

func (c *cart) GetList(param entity.CartParam) ([]entity.Cart, error) {
	carts := []entity.Cart{}
	if err := c.db.Where(param).Find(&carts).Error; err != nil {
		return carts, err
	}

	return carts, nil
}

func (c *cart) Get(param entity.CartParam) (entity.Cart, error) {
	cart := entity.Cart{}
	if err := c.db.Where(param).First(&cart).Error; err != nil {
		return cart, err
	}

	return cart, nil
}

func (c *cart) Update(selectParam entity.CartParam, updateParam entity.UpdateCartParam) error {
	if err := c.db.Model(entity.Cart{}).Where(selectParam).Updates(updateParam).Error; err != nil {
		return err
	}

	return nil
}
