package cache

import (
	"context"
)

func (c *cache) IncrUnreadDialogsCounter(ctx context.Context, userID string) (newCount int, err error) {
	val, err := c.redis.Incr(ctx, KeyUsersUnread+userID).Result()

	return int(val), err
}
