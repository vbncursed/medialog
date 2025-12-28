package main

import (
	"fmt"
	"os"

	"github.com/vbncursed/medialog/auth_service/config"
	"github.com/vbncursed/medialog/auth_service/internal/bootstrap"
)

func main() {
	bootstrap.InitLogger()

	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}

	redisClient := bootstrap.InitRedis(cfg)
	authStorage := bootstrap.InitPGStorage(cfg)
	sessionStorage := bootstrap.InitSessionStorage(redisClient)
	authService := bootstrap.InitAuthService(authStorage, sessionStorage, cfg)
	loginLimiter, registerLimiter, refreshLimiter := bootstrap.InitAuthRateLimiters(redisClient, cfg)
	authAPI := bootstrap.InitAuthServiceAPI(authService, loginLimiter, registerLimiter, refreshLimiter)

	bootstrap.AppRun(*authAPI)
}
