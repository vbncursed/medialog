package metadata_cache

import (
	"github.com/redis/go-redis/v9"
)

type MetadataCache struct {
	rdb *redis.Client
}

func NewMetadataCache(rdb *redis.Client) *MetadataCache {
	return &MetadataCache{
		rdb: rdb,
	}
}
