package cache

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func (c *cache) GetUnreadDialogsCounter(ctx context.Context, userID string) (count int, err error) {
	val, err := c.redis.Get(ctx, KeyUsersUnread+userID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}

		return 0, err
	}

	return strconv.Atoi(val)
}
