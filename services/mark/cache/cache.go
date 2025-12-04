package cache

import (
	"errors"
	"github.com/osamikoyo/music-and-marks/logger"
	"github.com/osamikoyo/music-and-marks/services/mark/config"
	"github.com/osamikoyo/music-and-marks/services/mark/entity"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

var (
	ErrCache       = errors.New("failed add value to cache")
	ErrConvertFail = errors.New("failed convert value")
)

type Cache struct {
	cache  *cache.Cache
	logger *logger.Logger
}

func NewCache(cfg *config.Config, logger *logger.Logger) *Cache {
	cache := cache.New(cfg.Cache.ExpTime, cfg.Cache.ExpiredItemsPurgeTimeout)

	return &Cache{
		cache:  cache,
		logger: logger,
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.logger.Info("setting value",
		zap.String("key", key),
		zap.Any("value", value))

	c.cache.Set(key, value, cache.DefaultExpiration)
}

func (c *Cache) GetReviews(key string) ([]entity.Review, error) {

	c.logger.Info("fetching reviews from cache",
		zap.String("key", key))

	value, ok := c.cache.Get(key)
	if !ok {
		c.logger.Error("failed fetch reviews from cache",
			zap.String("key", key))

		return nil, ErrCache
	}

	reviews, ok := value.([]entity.Review)
	if !ok {
		c.logger.Error("failed convert cache value to reviews",
			zap.String("key", key),
			zap.Any("value", value))

		return nil, ErrConvertFail
	}

	c.logger.Info("reviews fetched from cache successfully",
		zap.Any("reviews", reviews))

	return reviews, nil
}
