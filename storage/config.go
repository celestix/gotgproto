package storage

import (
	"fmt"
	"os"

	env "github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v2"
)

var DefaultStorageConfig = &StorageConfig{
	Cache: CacheConfig{
		Enabled:       true,
		InMemoryOnly:  false,
		TimeToLive:    "6h",
		CleanInterval: "24h",
		Revaluate:     true,
		InitialSize:   1000,
		ShardCount:    32,
	},
	Database: DatabaseConfig{
		MaxOpenConns:    100,
		MaxIdleConns:    10,
		ConnMaxLifetime: "1h",
		ConnMaxIdleTime: "30m",
		LogLevel:        "silent",
	},
}

type DatabaseConfig struct {
	MaxOpenConns    int    `yaml:"maxOpenConns" env:"DB_MAX_OPEN_CONNS" envDefault:"100"`
	MaxIdleConns    int    `yaml:"maxIdleConns" env:"DB_MAX_IDLE_CONNS" envDefault:"10"`
	ConnMaxLifetime string `yaml:"connMaxLifetime" env:"DB_CONN_MAX_LIFETIME" envDefault:"1h"`
	ConnMaxIdleTime string `yaml:"connMaxIdleTime" env:"DB_CONN_MAX_IDLE_TIME" envDefault:"30m"`
	LogLevel        string `yaml:"logLevel" env:"DB_LOG_LEVEL" envDefault:"silent"`
	TablePrefix     string `yaml:"tablePrefix" env:"DB_TABLE_PREFIX" envDefault:""`
}

type CacheConfig struct {
	Enabled       bool   `yaml:"enabled" env:"CACHE_ENABLED" envDefault:"true"`
	InMemoryOnly  bool   `yaml:"inMemoryOnly" env:"CACHE_IN_MEMORY_ONLY" envDefault:"false"`
	TimeToLive    string `yaml:"timeToLive" env:"CACHE_TTL" envDefault:"6h"`
	CleanInterval string `yaml:"cleanInterval" env:"CACHE_CLEAN_INTERVAL" envDefault:"24h"`
	Revaluate     bool   `yaml:"revaluate" env:"CACHE_REVALUATE" envDefault:"true"`
	InitialSize   int    `yaml:"initialSize" env:"CACHE_INITIAL_SIZE" envDefault:"1000"`
	ShardCount    int    `yaml:"shardCount" env:"CACHE_SHARD_COUNT" envDefault:"32"`
}

type StorageConfig struct {
	Database DatabaseConfig `yaml:"database"`
	Cache    CacheConfig    `yaml:"cache"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*StorageConfig, error) {
	cfg := &StorageConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// LoadConfigFromFile loads configuration from a YAML file and environment variables
func LoadConfigFromFile(path string) (*StorageConfig, error) {
	cfg := &StorageConfig{}

	// Read YAML file if provided
	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parsing config file: %w", err)
		}
	}

	// Environment variables override YAML settings
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parsing environment variables: %w", err)
	}

	return cfg, nil
}
