package cache

import (
	"context"
)

func (c *cache) DecrUnreadDialogsCounter(ctx context.Context, userID string) (newCount int, err error) {
	v, err := c.GetUnreadDialogsCounter(ctx, userID)
	if err == nil && v == 0 {
		return int(v), nil
	}

	val, err := c.redis.Decr(ctx, KeyUsersUnread+userID).Result()

	return int(val), err
}
