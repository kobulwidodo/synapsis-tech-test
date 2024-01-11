package category

import (
	"context"
	"encoding/json"
	"fmt"
	"go-clean/src/business/entity"
	"time"
)

const (
	getCategoryList = `synapsis:category:get`
)

func (c *cateogry) getCacheList(ctx context.Context) ([]entity.Category, error) {
	result := []entity.Category{}

	categoriesRedis, err := c.redis.Get(ctx, getCategoryList)
	if err != nil {
		return result, err
	}

	data := []byte(categoriesRedis)
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal redis : %v", err.Error())
	}

	return result, nil
}

func (c *cateogry) upsertCacheList(ctx context.Context, value []entity.Category, expTime time.Duration) error {
	categories, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := c.redis.SetEX(ctx, getCategoryList, string(categories), expTime); err != nil {
		return err
	}

	return nil
}
