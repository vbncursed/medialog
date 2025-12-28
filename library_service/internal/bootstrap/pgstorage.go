package bootstrap

import (
	"fmt"
	"log/slog"

	"github.com/vbncursed/medialog/library_service/config"
	"github.com/vbncursed/medialog/library_service/internal/storage/library_storage"
)

func InitPGStorage(cfg *config.Config) *library_storage.LibraryStorage {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName, cfg.Database.SSLMode)
	storage, err := library_storage.NewLibraryStorage(connectionString)
	if err != nil {
		slog.Error("ошибка инициализации БД", "err", err)
		panic(err)
	}
	return storage
}

