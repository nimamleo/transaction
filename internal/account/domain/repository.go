package domain

import "context"

type AccountRepository interface {
	Create(ctx context.Context, account *Account) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Account, error)
	GetByUserID(ctx context.Context, userID string) ([]*Account, error)
	UpdateBalance(ctx context.Context, id string, balance int64) error

	CreateSystemAccount(ctx context.Context, systemAccount *SystemAccount) error
	GetSystemAccountByCurrency(ctx context.Context, currency Currency) (*SystemAccount, error)
	SystemAccountExistsByCurrency(ctx context.Context, currency Currency) (bool, error)

	CreateTransaction(ctx context.Context, transaction *Transaction) error
	GetTransactionByReference(ctx context.Context, reference string) (*Transaction, error)
	TransactionExistsByReference(ctx context.Context, reference string) (bool, error)
}
