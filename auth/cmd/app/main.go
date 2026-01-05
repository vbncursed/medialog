package main

import (
	"fmt"
	"os"

	"github.com/vbncursed/medialog/auth/config"
	"github.com/vbncursed/medialog/auth/internal/bootstrap"
)

func main() {
	configPath := os.Getenv("configPath")
	if configPath == "" {
		configPath = "./config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	authStorage := bootstrap.InitPGStorage(cfg)
	redisClient := bootstrap.InitRedis(cfg)
	sessionStorage := bootstrap.InitSessionStorage(redisClient)

	authService := bootstrap.InitAuthService(authStorage, sessionStorage, cfg)

	loginLimiter, registerLimiter, refreshLimiter := bootstrap.InitRateLimiters(redisClient, cfg)

	authAPI := bootstrap.InitAuthServiceAPI(authService, cfg, loginLimiter, registerLimiter, refreshLimiter)

	bootstrap.AppRun(authAPI, cfg)
}
