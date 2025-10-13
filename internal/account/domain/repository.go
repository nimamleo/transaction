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

	GetTransactionByReference(ctx context.Context, reference string) (*Transaction, error)
	TransactionExistsByReference(ctx context.Context, reference string, accountID string) (bool, error)
	CreateTransactionAndUpdateBalance(ctx context.Context, transaction *Transaction, accountID string, newBalance int64) (*Transaction, error)
	CreateTransferTransactions(ctx context.Context, fromAccountID, toAccountID, reference string, amount int64, fromNewBalance, toNewBalance int64) error
	GetAccountTransactions(ctx context.Context, accountID string, limit int, after string) ([]*Transaction, error)
}
