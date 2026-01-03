package metadata_cache

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (c *MetadataCache) GetMedia(ctx context.Context, key string) (*models.Media, error) {
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

	var media models.Media
	if err := json.Unmarshal([]byte(data), &media); err != nil {
		return nil, err
	}

	return &media, nil
}
