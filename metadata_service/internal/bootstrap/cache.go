package bootstrap

import (
	"github.com/redis/go-redis/v9"
	"github.com/vbncursed/medialog/metadata_service/internal/storage/metadata_cache"
)

func InitMetadataCache(redisClient *redis.Client) *metadata_cache.MetadataCache {
	return metadata_cache.NewMetadataCache(redisClient)
}

