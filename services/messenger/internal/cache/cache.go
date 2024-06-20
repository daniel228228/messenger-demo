package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"messenger.messenger/pkg/log"

	"messenger.messenger/pkg/config"
)

var Name = "Cache"

type Cache interface {
	IncrUnreadDialogsCounter(ctx context.Context, userID string) (newCount int, err error)
	DecrUnreadDialogsCounter(ctx context.Context, userID string) (newCount int, err error)
	GetUnreadDialogsCounter(ctx context.Context, userID string) (count int, err error)
}

const (
	ConnRetries = 20
)

var (
	KeyUsersUnread = "users:unread:"
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

func (c *cache) connect(ctx context.Context) error {
	redis := redis.NewClient(&redis.Options{
		Addr:     c.config.CacheUrl,
		Password: "",
		DB:       0,
	})

	ping, err := redis.Ping(ctx).Result()
	if err != nil {
		return err
	}

	c.log.Trace().Str("from", "cache").Msgf("Successful Redis connection: %s", ping)
	c.redis = redis

	return nil
}

func (c *cache) Stop(_ context.Context) error {
	c.log.Info().Bool("app", true).Str("component", Name).Str("state", "stop").Send()

	return c.redis.Close()
}
