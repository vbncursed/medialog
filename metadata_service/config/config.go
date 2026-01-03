package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	Database     DatabaseConfig     `yaml:"database"`
	Redis        RedisConfig        `yaml:"redis"`
	Kafka        KafkaConfig        `yaml:"kafka"`
	Server       ServerConfig       `yaml:"server"`
	Auth         AuthConfig         `yaml:"auth"`
	ExternalAPIs ExternalAPIsConfig `yaml:"external_apis"`
	Cache        CacheConfig        `yaml:"cache"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type KafkaConfig struct {
	Host                     string `yaml:"host"`
	Port                     int    `yaml:"port"`
	MetadataEventTopic       string `yaml:"metadata_event_topic"`
	LibraryEntryChangedTopic string `yaml:"library_entry_changed_topic"`
}

type ServerConfig struct {
	GRPCAddr string `yaml:"grpc_addr"`
	HTTPAddr string `yaml:"http_addr"`
}

type AuthConfig struct {
	JWTSecret string `yaml:"jwt_secret"`
}

type ExternalAPIsConfig struct {
	TMDB        TMDBConfig        `yaml:"tmdb"`
	OpenLibrary OpenLibraryConfig `yaml:"open_library"`
}

type TMDBConfig struct {
	APIKey         string `yaml:"api_key"`
	BaseURL        string `yaml:"base_url"`
	ImageBaseURL   string `yaml:"image_base_url"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
}

type OpenLibraryConfig struct {
	BaseURL        string `yaml:"base_url"`
	CoverBaseURL   string `yaml:"cover_base_url"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
}

type CacheConfig struct {
	MediaTTLSeconds  int64 `yaml:"media_ttl_seconds"`  // TTL для кэша метаданных
	SearchTTLSeconds int64 `yaml:"search_ttl_seconds"` // TTL для кэша результатов поиска
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &cfg, nil
}
