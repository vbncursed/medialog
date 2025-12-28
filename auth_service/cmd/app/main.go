package main

import (
	"fmt"
	"os"
	"time"

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
	loginLimiter, registerLimiter, refreshLimiter := bootstrap.InitAuthRateLimiters(redisClient, cfg.Auth.RateLimitLoginPerMinute, cfg.Auth.RateLimitRegisterPerMinute, cfg.Auth.RateLimitRefreshPerMinute)
	authAPI := bootstrap.InitAuthServiceAPI(authService, loginLimiter, registerLimiter, refreshLimiter)

	// Запускаем периодическую очистку старых сессий
	cleanupInterval := time.Duration(cfg.Auth.SessionCleanupIntervalHours) * time.Hour
	retentionPeriod := time.Duration(cfg.Auth.SessionRetentionPeriodDays) * 24 * time.Hour
	if cleanupInterval > 0 && retentionPeriod > 0 {
		bootstrap.StartSessionCleanup(authStorage, cleanupInterval, retentionPeriod)
	}

	bootstrap.AppRun(*authAPI)
}
