package bootstrap

import (
	server "github.com/vbncursed/medialog/metadata_service/internal/api/metadata_service_api"
	"github.com/vbncursed/medialog/metadata_service/config"
	"github.com/vbncursed/medialog/metadata_service/internal/services/metadata_service"
)

func InitMetadataServiceAPI(metadataService *metadata_service.MetadataService, cfg *config.Config) *server.MetadataServiceAPI {
	return server.NewMetadataServiceAPI(metadataService, cfg.Auth.JWTSecret)
}

