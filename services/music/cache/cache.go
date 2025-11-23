// Package cache stores cacher
package cache

import (
	"errors"
	"fmt"

	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/music/config"
	"github.com/osamikoyo/music-and-marks/services/music/entity"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

var (
	ErrNilInput      = errors.New("empty fields")
	ErrConvertFailed = errors.New("failed convert cache value to entity")
)

type Cache struct {
	logger *logger.Logger
	cache  *cache.Cache
}

func NewCache(cfg *config.Config, logger *logger.Logger) *Cache {
	return &Cache{
		logger: logger,
		cache:  cache.New(cfg.Cache.ExpTime, cfg.Cache.ExpiredItemsPurgeTimeout),
	}
}

func (c *Cache) SetArtist(key string, artist *entity.Artist) error {
	if artist == nil || key == "" {
		return ErrNilInput
	}

	c.logger.Info("setting artist in cache",
		zap.Any("artist", artist),
		zap.String("key", key))

	if err := c.cache.Add(key, artist, cache.DefaultExpiration); err != nil {
		c.logger.Error("failed set artist",
			zap.String("key", key),
			zap.Error(err))

		return fmt.Errorf("failed set artist: %w", err)
	}

	c.logger.Info("artist set in cache",
		zap.Any("artist", artist),
		zap.String("key", key))

	return nil
}

func (c *Cache) SetRelease(key string, release *entity.Release) error {
	if release == nil || key == "" {
		return ErrNilInput
	}

	c.logger.Info("setting release in cache",
		zap.Any("release", release),
		zap.String("key", key))

	if err := c.cache.Add(key, release, cache.DefaultExpiration); err != nil {
		c.logger.Error("failed set release",
			zap.String("key", key),
			zap.Error(err))

		return fmt.Errorf("failed set releases: %w", err)
	}

	c.logger.Info("release set in cache",
		zap.Any("release", release),
		zap.String("key", key))

	return nil
}

func (c *Cache) GetArtist(key string) (*entity.Artist, error) {
	if key == "" {
		return nil, ErrNilInput
	}

	c.logger.Info("fetching artist",
		zap.String("key", key))

	value, ok := c.cache.Get(key)
	if !ok {
		c.logger.Error("failed fetch artist from cache",
			zap.String("key", key))

		return nil, fmt.Errorf("failed fetch artist from cache with key: %s", key)
	}

	artist, ok := value.(entity.Artist)
	if !ok {
		c.logger.Error("failed convert cache value to entity",
			zap.Any("value", value))

		return nil, ErrConvertFailed
	}

	return &artist, nil
}

func (c *Cache) GetRelease(key string) (*entity.Release, error) {
	if key == "" {
		return nil, ErrNilInput
	}

	c.logger.Info("fetching release",
		zap.String("key", key))

	value, ok := c.cache.Get(key)
	if !ok {
		c.logger.Error("failed fetch release from cache",
			zap.String(key, key))

		return nil, fmt.Errorf("failed fetch artist from cache with key: %s", key)
	}

	release, ok := value.(entity.Release)
	if !ok {
		c.logger.Error("failed convert cache value to entity",
			zap.Any("value", value))

		return nil, ErrConvertFailed
	}

	return &release, nil
}
