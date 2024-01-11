package category

import (
	"context"
	"errors"
	"go-clean/src/business/entity"
	"go-clean/src/lib/redis"
	"log"
	"time"

	"gorm.io/gorm"
)

type Interface interface {
	GetList(ctx context.Context) ([]entity.Category, error)
}

type cateogry struct {
	db    *gorm.DB
	redis redis.Interface
}

func Init(db *gorm.DB, redis redis.Interface) Interface {
	c := &cateogry{
		db:    db,
		redis: redis,
	}

	return c
}

func (c *cateogry) GetList(ctx context.Context) ([]entity.Category, error) {
	cacheResult, err := c.getCacheList(ctx)
	switch {
	case errors.Is(err, redis.Nil):
		log.Printf("error redis is nil, %s\n", err.Error())
	case err != nil:
		log.Printf("error redis : %s\n", err.Error())
	default:
		return cacheResult, nil
	}

	categories := []entity.Category{}
	if err := c.db.Find(&categories).Error; err != nil {
		return categories, err
	}

	if err := c.upsertCacheList(ctx, categories, time.Minute); err != nil {
		log.Println(err.Error())
	}

	return categories, nil
}
