package metadata_cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (c *MetadataCache) SetSearchResults(ctx context.Context, key string, results *models.SearchMediaResult, ttl int64) error {
	if c.rdb == nil || results == nil {
		return nil
	}

	data, err := json.Marshal(results)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err()
}

