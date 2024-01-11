package product

import (
	"context"
	"encoding/json"
	"fmt"
	"go-clean/src/business/entity"
	"time"
)

const (
	getProductList    = `synapsis:product:get:q:%s`
	getProductByIdKey = `synapsis:product:get:%s`
)

func (p *product) getCacheList(ctx context.Context, marshalledParams []byte) ([]entity.Product, error) {
	result := []entity.Product{}

	productListKey := fmt.Sprintf(getProductList, marshalledParams)

	categoriesRedis, err := p.redis.Get(ctx, productListKey)
	if err != nil {
		return result, err
	}

	data := []byte(categoriesRedis)
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal redis : %v", err.Error())
	}

	return result, nil
}

func (p *product) upsertCacheList(ctx context.Context, key []byte, value []entity.Product, expTime time.Duration) error {
	categories, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := p.redis.SetEX(ctx, fmt.Sprintf(getProductList, string(key)), string(categories), expTime); err != nil {
		return err
	}

	return nil
}

func (p *product) getCacheByID(ctx context.Context, marshalledParams []byte) (entity.Product, error) {
	var result entity.Product
	parameterByIDKey := fmt.Sprintf(getProductByIdKey, marshalledParams)
	parameter, err := p.redis.Get(ctx, parameterByIDKey)
	if err != nil {
		return result, err
	}

	data := []byte(parameter)
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal redis : %v", err.Error())
	}

	return result, nil
}

func (p *product) upsertCacheByID(ctx context.Context, marshalledParams []byte, product entity.Product, expTime time.Duration) error {
	parameterByIDKey := fmt.Sprintf(getProductByIdKey, marshalledParams)

	rawJSON, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal redis : %v", err.Error())
	}

	return p.redis.SetEX(ctx, parameterByIDKey, string(rawJSON), expTime)
}
