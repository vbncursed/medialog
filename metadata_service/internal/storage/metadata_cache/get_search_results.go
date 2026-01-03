package metadata_cache

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (c *MetadataCache) GetSearchResults(ctx context.Context, key string) (*models.SearchMediaResult, error) {
	if c.rdb == nil {
		return nil, nil
	}

	data, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var result models.SearchMediaResult
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

