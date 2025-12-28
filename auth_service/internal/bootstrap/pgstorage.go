package bootstrap

import (
	"fmt"
	"log/slog"

	"github.com/vbncursed/medialog/auth_service/config"
	"github.com/vbncursed/medialog/auth_service/internal/storage/auth_storage"
)

func InitPGStorage(cfg *config.Config) *auth_storage.AuthStorage {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	storage, err := auth_storage.NewAuthStorage(connectionString)
	if err != nil {
		slog.Error("ошибка инициализации БД", "err", err)
		panic(err)
	}
	return storage
}
