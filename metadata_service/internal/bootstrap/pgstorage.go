package bootstrap

import (
	"fmt"
	"log/slog"

	"github.com/vbncursed/medialog/metadata_service/config"
	"github.com/vbncursed/medialog/metadata_service/internal/storage/metadata_storage"
)

func InitPGStorage(cfg *config.Config) *metadata_storage.MetadataStorage {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName, cfg.Database.SSLMode)
	storage, err := metadata_storage.NewMetadataStorage(connectionString)
	if err != nil {
		slog.Error("ошибка инициализации БД", "err", err)
		panic(err)
	}
	return storage
}

