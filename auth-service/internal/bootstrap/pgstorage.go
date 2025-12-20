package bootstrap

import (
	"fmt"
	"log/slog"

	"github.com/vbncursed/medialog/auth-service/config"
	"github.com/vbncursed/medialog/auth-service/internal/storage/pgstorage"
)

func InitPGStorage(cfg *config.Config) *pgstorage.PGstorage {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	storage, err := pgstorage.NewPGStorge(connectionString)
	if err != nil {
		slog.Error("ошибка инициализации БД", "err", err)
		panic(err)
	}
	return storage
}
