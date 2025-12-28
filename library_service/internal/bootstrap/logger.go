package bootstrap

import (
	"log/slog"
	"os"
	"strings"
	"time"
)

// InitLogger настраивает slog по env:
// - LOG_FORMAT: "json" | "text" (default: text)
// - LOG_LEVEL: "debug" | "info" | "warn" | "error" (default: info)
func InitLogger() {
	level := parseLevel(os.Getenv("LOG_LEVEL"))
	format := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_FORMAT")))
	if format == "" {
		// Для локальной разработки читаемее text.
		format = "text"
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: os.Getenv("LOG_ADD_SOURCE") == "1",
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				// Короткое время для консоли.
				if t, ok := a.Value.Any().(time.Time); ok {
					return slog.String(slog.TimeKey, t.Format("15:04:05"))
				}
			case slog.LevelKey:
				// INFO/WARN/ERROR/DEBUG
				if lv, ok := a.Value.Any().(slog.Level); ok {
					return slog.String(slog.LevelKey, strings.ToUpper(lv.String()))
				}
			}
			return a
		},
	}

	var h slog.Handler
	switch format {
	case "text":
		h = slog.NewTextHandler(os.Stdout, opts)
	default:
		// Для JSON обычно полезно видеть source.
		opts.AddSource = true
		h = slog.NewJSONHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(h))
}

func parseLevel(v string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

