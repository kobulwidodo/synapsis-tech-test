package product

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	GetList(param entity.ProductParam) ([]entity.Product, error)
	GetListByID(productIDs []uint) ([]entity.Product, error)
	Get(param entity.ProductParam) (entity.Product, error)
}

type product struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	p := &product{
		db: db,
	}

	return p
}

func (p *product) GetList(param entity.ProductParam) ([]entity.Product, error) {
	products := []entity.Product{}
	if err := p.db.Where(param).Find(&products).Error; err != nil {
		return products, err
	}
	return products, nil
}

func (p *product) GetListByID(productIDs []uint) ([]entity.Product, error) {
	products := []entity.Product{}
	if err := p.db.Find(&products, productIDs).Error; err != nil {
		return products, err
	}

	return products, nil
}

func (p *product) Get(param entity.ProductParam) (entity.Product, error) {
	product := entity.Product{}
	if err := p.db.Where(param).First(&product).Error; err != nil {
		return product, err
	}
	return product, nil
}
