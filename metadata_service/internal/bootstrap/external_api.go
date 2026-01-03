package bootstrap

import (
	"log/slog"

	"github.com/vbncursed/medialog/metadata_service/config"
	"github.com/vbncursed/medialog/metadata_service/internal/clients"
	"github.com/vbncursed/medialog/metadata_service/internal/clients/open_library"
	"github.com/vbncursed/medialog/metadata_service/internal/clients/tmdb"
)

func InitExternalAPIClient(cfg *config.Config) *clients.ExternalAPIClient {
	var tmdbClient *tmdb.TMDBClient
	if cfg.ExternalAPIs.TMDB.APIKey != "" && cfg.ExternalAPIs.TMDB.APIKey != "CHANGE_ME" {
		slog.Info("Initializing TMDB client", "baseURL", cfg.ExternalAPIs.TMDB.BaseURL)
		tmdbClient = tmdb.NewTMDBClient(
			cfg.ExternalAPIs.TMDB.APIKey,
			cfg.ExternalAPIs.TMDB.BaseURL,
			cfg.ExternalAPIs.TMDB.ImageBaseURL,
			cfg.ExternalAPIs.TMDB.TimeoutSeconds,
		)
	} else {
		slog.Warn("TMDB client not initialized: API key not configured or is CHANGE_ME")
	}

	slog.Info("Initializing OpenLibrary client", "baseURL", cfg.ExternalAPIs.OpenLibrary.BaseURL)
	openLibClient := open_library.NewOpenLibraryClient(
		cfg.ExternalAPIs.OpenLibrary.BaseURL,
		cfg.ExternalAPIs.OpenLibrary.CoverBaseURL,
		cfg.ExternalAPIs.OpenLibrary.TimeoutSeconds,
	)

	return clients.NewExternalAPIClient(tmdbClient, openLibClient)
}
