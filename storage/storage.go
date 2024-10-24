package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AnimeKaizoku/cacher"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type PeerStorage struct {
	peerCache  *cacher.Cacher[int64, *Peer]
	peerLock   *sync.RWMutex
	config     *StorageConfig
	SqlSession *gorm.DB
}

type Option func(*PeerStorage) error

func WithCustomCache(cache *cacher.Cacher[int64, *Peer]) Option {
	return func(ps *PeerStorage) error {
		ps.peerCache = cache
		return nil
	}
}

func NewPeerStorage(ctx context.Context, cfg *StorageConfig, dialector gorm.Dialector, opts ...Option) (*PeerStorage, error) {
	if cfg == nil {
		return nil, ErrInvalidConfig
	}

	ps := &PeerStorage{
		config:   cfg,
		peerLock: new(sync.RWMutex),
	}

	for _, opt := range opts {
		if err := opt(ps); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	if !cfg.Cache.Enabled || !cfg.Cache.InMemoryOnly {
		if err := ps.initDatabase(dialector); err != nil {
			return nil, fmt.Errorf("failed to initialize database: %w", err)
		}
	}

	if err := ps.initCache(); err != nil {
		return nil, fmt.Errorf("failed to initialize cache: %w", err)
	}

	return ps, nil
}

func (ps *PeerStorage) initDatabase(dialector gorm.Dialector) error {
	logLevel := logger.Silent
	switch ps.config.Database.LogLevel {
	case "info":
		logLevel = logger.Info
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	}

	config := gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logLevel),
	}

	if ps.config.Database.TablePrefix != "" {
		config.NamingStrategy = schema.NamingStrategy{
			TablePrefix: ps.config.Database.TablePrefix,
		}
	}

	db, err := gorm.Open(dialector, &config)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	maxLifetime, err := time.ParseDuration(ps.config.Database.ConnMaxLifetime)
	if err != nil {
		return fmt.Errorf("invalid connection max lifetime: %w", err)
	}

	maxIdleTime, err := time.ParseDuration(ps.config.Database.ConnMaxIdleTime)
	if err != nil {
		return fmt.Errorf("invalid connection max idle time: %w", err)
	}

	sqlDB.SetMaxOpenConns(ps.config.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(ps.config.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(maxLifetime)
	sqlDB.SetConnMaxIdleTime(maxIdleTime)

	if err := db.AutoMigrate(&Session{}, &Peer{}); err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	ps.SqlSession = db

	return nil
}

func (ps *PeerStorage) initCache() error {
	ttl, err := time.ParseDuration(ps.config.Cache.TimeToLive)
	if err != nil {
		return fmt.Errorf("invalid cache TTL: %w", err)
	}

	cleanInterval, err := time.ParseDuration(ps.config.Cache.CleanInterval)
	if err != nil {
		return fmt.Errorf("invalid cache clean interval: %w", err)
	}

	opts := &cacher.NewCacherOpts{
		TimeToLive:    ttl,
		CleanInterval: cleanInterval,
		Revaluate:     ps.config.Cache.Revaluate,
	}

	ps.peerCache = cacher.NewCacher[int64, *Peer](opts)

	return nil
}
