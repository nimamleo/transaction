package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"transaction/internal/account/domain"

	"github.com/redis/go-redis/v9"
)

type balanceCacheData struct {
	Balance   int64     `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

type accountCache struct {
	client *redis.Client
}

func NewAccountCache(client *redis.Client) domain.AccountCache {
	return &accountCache{client: client}
}

func (c *accountCache) GetBalance(ctx context.Context, accountID string) (*domain.BalanceCache, error) {
	key := fmt.Sprintf("account:balance:%s", accountID)

	data, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var balanceCache balanceCacheData
	if err := json.Unmarshal([]byte(data), &balanceCache); err != nil {
		return nil, err
	}

	return &domain.BalanceCache{
		Balance:   balanceCache.Balance,
		UpdatedAt: balanceCache.UpdatedAt,
	}, nil
}

func (c *accountCache) SetBalance(ctx context.Context, accountID string, balance int64, updatedAt time.Time) error {
	key := fmt.Sprintf("account:balance:%s", accountID)

	balanceCache := balanceCacheData{
		Balance:   balance,
		UpdatedAt: updatedAt,
	}

	data, err := json.Marshal(balanceCache)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, 1*time.Minute).Err()
}
