package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"

	"messenger.auth/pkg/log"

	"messenger.auth/pkg/config"
)

var Name = "Cache"

type Cache interface {
	Set(id, userID string, exp time.Duration) error
	Get(id string) (string, error)
	Del(id string) error
}

var (
	ErrEmpty = errors.New("empty result")
)

const (
	ConnRetries = 20
)

type cache struct {
	config *config.Config
	log    *log.Logger
	redis  *redis.Client
}

func NewCache(config *config.Config, log *log.Logger) *cache {
	return &cache{
		config: config,
		log:    log,
	}
}

func (c *cache) Start(ctx context.Context) error {
	period := 16 * time.Millisecond

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for i := 0; i < ConnRetries; i++ {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}

		c.log.Trace().Str("from", "cache").Msgf("Attempt %d: connecting to Redis", i+1)

		if err := c.connect(ctx); err != nil {
			c.log.Trace().Str("from", "cache").Msgf("Attempt %d failed: Redis error: %s", i+1, err.Error())

			if i+1 < ConnRetries {
				if period < 5*time.Second {
					period *= 2
				} else {
					period = 10 * time.Second
				}

				ticker.Reset(period)
			} else {
				c.log.Error().Err(err).Msg("Can't connect to Redis")
				return err
			}
		} else {
			c.log.Info().Bool("app", true).Str("component", Name).Str("state", "start").Send()
			return nil
		}
	}

	return nil
}

func (c *cache) connect(_ context.Context) error {
	redis := redis.NewClient(&redis.Options{
		Addr:     c.config.CacheUrl,
		Password: "",
		DB:       0,
	})

	ping, err := redis.Ping().Result()
	if err != nil {
		return err
	}

	c.log.Trace().Str("from", "cache").Msgf("Successful Redis connection: %s", ping)
	c.redis = redis

	return nil
}

func (c *cache) Stop(_ context.Context) error {
	c.log.Info().Bool("app", true).Str("component", Name).Str("state", "stop").Send()

	return nil
}

func (c *cache) Set(id, userID string, exp time.Duration) error {
	return c.redis.Set(id, userID, exp).Err()
}

func (c *cache) Get(id string) (string, error) {
	res, err := c.redis.Get(id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrEmpty
		}

		return "", err
	}

	return res, nil
}

func (c *cache) Del(id string) error {
	err := c.redis.Del(id).Err()

	if errors.Is(err, redis.Nil) {
		return ErrEmpty
	}

	return err
}
