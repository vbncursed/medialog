package bootstrap

import (
	"github.com/vbncursed/medialog/metadata_service/config"
	"github.com/vbncursed/medialog/metadata_service/internal/services/metadata_service"
	"github.com/vbncursed/medialog/metadata_service/internal/storage/metadata_cache"
	"github.com/vbncursed/medialog/metadata_service/internal/storage/metadata_storage"
)

func InitMetadataService(
	storage *metadata_storage.MetadataStorage,
	cache *metadata_cache.MetadataCache,
	externalAPI metadata_service.ExternalAPIClient,
	cfg *config.Config,
) *metadata_service.MetadataService {
	return metadata_service.NewMetadataService(
		storage,
		cache,
		externalAPI,
		cfg.Cache.MediaTTLSeconds,
		cfg.Cache.SearchTTLSeconds,
	)
}
