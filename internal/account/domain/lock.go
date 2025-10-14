package domain

import (
	"context"
	"time"
)

type Lock interface {
	Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error)
	Release(ctx context.Context, key string) error
	Extend(ctx context.Context, key string, ttl time.Duration) error
}
