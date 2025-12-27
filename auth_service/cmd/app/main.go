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

	redisClient := bootstrap.InitRedis(cfg)
	authStorage := bootstrap.InitPGStorage(cfg)
	authService := bootstrap.InitAuthService(authStorage, cfg)
	loginLimiter, registerLimiter := bootstrap.InitAuthRateLimiters(redisClient, cfg.Auth.RateLimitLoginPerMinute, cfg.Auth.RateLimitRegisterPerMinute)
	authAPI := bootstrap.InitAuthServiceAPI(authService, loginLimiter, registerLimiter)

	bootstrap.AppRun(*authAPI)
}
