package metadata_cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/vbncursed/medialog/metadata_service/internal/models"
)

func (c *MetadataCache) SetMedia(ctx context.Context, key string, media *models.Media, ttl int64) error {
	if c.rdb == nil || media == nil {
		return nil
	}

	data, err := json.Marshal(media)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err()
}

