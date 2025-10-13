package domain

import "context"

type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Account, error)
	GetByUserID(ctx context.Context, userID string) ([]*Account, error)
	UpdateBalance(ctx context.Context, id string, balance int64) error
}
