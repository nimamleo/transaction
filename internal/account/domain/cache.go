package domain

import (
	"context"
	"time"
)

type BalanceCache struct {
	Balance   int64
	UpdatedAt time.Time
}

type AccountCache interface {
	GetBalance(ctx context.Context, accountID string) (*BalanceCache, error)
	SetBalance(ctx context.Context, accountID string, balance int64, updatedAt time.Time) error
}
