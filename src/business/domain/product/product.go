package product

import (
	"context"
	"encoding/json"
	"errors"
	"go-clean/src/business/entity"
	"go-clean/src/lib/redis"
	"log"
	"time"

	"gorm.io/gorm"
)

type Interface interface {
	GetList(ctx context.Context, param entity.ProductParam) ([]entity.Product, error)
	GetListByID(ctx context.Context, productIDs []uint) ([]entity.Product, error)
	Get(ctx context.Context, param entity.ProductParam) (entity.Product, error)
}

type product struct {
	db    *gorm.DB
	redis redis.Interface
}

func Init(db *gorm.DB, redis redis.Interface) Interface {
	p := &product{
		db:    db,
		redis: redis,
	}

	return p
}

func (p *product) GetList(ctx context.Context, param entity.ProductParam) ([]entity.Product, error) {
	marshalledParam, err := json.Marshal(param)
	if err != nil {
		log.Println(err.Error())
	}

	cacheResult, err := p.getCacheList(ctx, marshalledParam)
	switch {
	case errors.Is(err, redis.Nil):
		log.Printf("error redis is nil, %s\n", err.Error())
	case err != nil:
		log.Printf("error redis : %s\n", err.Error())
	default:
		return cacheResult, nil
	}

	products := []entity.Product{}
	if err := p.db.Where(param).Find(&products).Error; err != nil {
		return products, err
	}

	key, err := json.Marshal(param)
	if err != nil {
		log.Println(err.Error())
	}

	if err := p.upsertCacheList(ctx, key, products, time.Minute); err != nil {
		log.Println(err.Error())
	}

	return products, nil
}

func (p *product) GetListByID(ctx context.Context, productIDs []uint) ([]entity.Product, error) {
	marshalledParam, err := json.Marshal(productIDs)
	if err != nil {
		log.Println(err.Error())
	}

	cacheResult, err := p.getCacheList(ctx, marshalledParam)
	switch {
	case errors.Is(err, redis.Nil):
		log.Printf("error redis is nil, %s\n", err.Error())
	case err != nil:
		log.Printf("error redis : %s\n", err.Error())
	default:
		return cacheResult, nil
	}

	products := []entity.Product{}
	if err := p.db.Find(&products, productIDs).Error; err != nil {
		return products, err
	}

	key, err := json.Marshal(productIDs)
	if err != nil {
		log.Println(err.Error())
	}

	if err := p.upsertCacheList(ctx, key, products, time.Minute); err != nil {
		log.Println(err.Error())
	}

	return products, nil
}

func (p *product) Get(ctx context.Context, param entity.ProductParam) (entity.Product, error) {
	marshalledParam, err := json.Marshal(param)
	if err != nil {
		log.Println(err.Error())
	}

	cacheResult, err := p.getCacheByID(ctx, marshalledParam)
	switch {
	case errors.Is(err, redis.Nil):
		log.Printf("error redis is nil, %s\n", err.Error())
	case err != nil:
		log.Printf("error redis : %s\n", err.Error())
	default:
		return cacheResult, nil
	}

	product := entity.Product{}
	if err := p.db.Where(param).First(&product).Error; err != nil {
		return product, err
	}

	if err = p.upsertCacheByID(ctx, marshalledParam, product, time.Minute); err != nil {
		log.Println(err.Error())
	}

	return product, nil
}
