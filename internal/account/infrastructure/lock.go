package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Lock struct {
	redisClient *redis.Client
}

func NewLock(redisClient *redis.Client) *Lock {
	return &Lock{
		redisClient: redisClient,
	}
}

func (l *Lock) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("lock:%s", key)

	acquired, err := l.redisClient.SetNX(ctx, lockKey, "locked", ttl).Result()
	if err != nil {
		return false, err
	}

	return acquired, nil
}

func (l *Lock) Release(ctx context.Context, key string) error {
	lockKey := fmt.Sprintf("lock:%s", key)
	return l.redisClient.Del(ctx, lockKey).Err()
}

func (l *Lock) Extend(ctx context.Context, key string, ttl time.Duration) error {
	lockKey := fmt.Sprintf("lock:%s", key)
	return l.redisClient.Expire(ctx, lockKey, ttl).Err()
}
