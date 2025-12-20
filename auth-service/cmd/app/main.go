package main

import (
	"fmt"
	"os"

	"github.com/vbncursed/medialog/auth-service/config"
	"github.com/vbncursed/medialog/auth-service/internal/bootstrap"
)

func main() {
	bootstrap.InitLogger()

	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}

	authStorage := bootstrap.InitPGStorage(cfg)
	authService := bootstrap.InitAuthService(authStorage, cfg)
	authAPI := bootstrap.InitAuthServiceAPI(authService, cfg)

	bootstrap.AppRun(*authAPI, cfg.Server.GRPCAddr, cfg.Server.HTTPAddr)
}
