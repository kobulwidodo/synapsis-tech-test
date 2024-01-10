package transaction

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(transaction entity.Transaction) (entity.Transaction, error)
}

type transaction struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	t := &transaction{
		db: db,
	}

	return t
}

func (t *transaction) Create(transaction entity.Transaction) (entity.Transaction, error) {
	if err := t.db.Create(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}
