package main

import (
	"fmt"
	"os"

	"github.com/vbncursed/medialog/library_service/config"
	"github.com/vbncursed/medialog/library_service/internal/bootstrap"
)

func main() {
	bootstrap.InitLogger()

	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}

	storage := bootstrap.InitPGStorage(cfg)
	producer := bootstrap.InitLibraryEntryEventProducer(cfg)
	libraryService := bootstrap.InitLibraryService(storage, producer)
	libraryAPI := bootstrap.InitLibraryServiceAPI(libraryService, cfg)

	bootstrap.AppRun(*libraryAPI, cfg)
}
