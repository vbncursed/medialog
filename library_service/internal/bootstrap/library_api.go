package bootstrap

import (
	"github.com/vbncursed/medialog/library_service/config"
	server "github.com/vbncursed/medialog/library_service/internal/api/library_service_api"
	"github.com/vbncursed/medialog/library_service/internal/services/library_service"
)

func InitLibraryServiceAPI(libraryService *library_service.LibraryService, cfg *config.Config) *server.LibraryServiceAPI {
	return server.NewLibraryServiceAPI(libraryService, cfg.Auth.JWTSecret)
}
