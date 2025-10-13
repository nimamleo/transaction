package domain

import "context"

type Ledger interface {
	CreateAccount(ctx context.Context, currency Currency) (string, error)
	GetBalance(ctx context.Context, ledgerID string) (int64, error)
}
