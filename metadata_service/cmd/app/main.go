package main

import (
	"fmt"
	"os"

	"github.com/vbncursed/medialog/metadata_service/config"
	"github.com/vbncursed/medialog/metadata_service/internal/bootstrap"
)

func main() {
	bootstrap.InitLogger()

	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}

	storage := bootstrap.InitPGStorage(cfg)
	redisClient := bootstrap.InitRedis(cfg)
	cache := bootstrap.InitMetadataCache(redisClient)
	externalAPI := bootstrap.InitExternalAPIClient(cfg)

	metadataService := bootstrap.InitMetadataService(
		storage,
		cache,
		externalAPI,
		cfg,
	)

	libraryEntryProcessor := bootstrap.InitLibraryEntryProcessor(metadataService)
	libraryEntryChangedConsumer := bootstrap.InitLibraryEntryChangedConsumer(cfg, libraryEntryProcessor)
	metadataAPI := bootstrap.InitMetadataServiceAPI(metadataService, cfg)

	bootstrap.AppRun(*metadataAPI, libraryEntryChangedConsumer, cfg)
}
