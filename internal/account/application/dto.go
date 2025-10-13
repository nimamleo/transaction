package application

import "time"

type BalanceInfo struct {
	Balance   int64
	UpdatedAt time.Time
}

type DepositResult struct {
	TransactionID string
	TransferID    string
	Amount        int64
	NewBalance    int64
	Status        string
}

type TransferResult struct {
	TransferID     string
	FromAccountID  string
	ToAccountID    string
	Amount         int64
	FromNewBalance int64
	ToNewBalance   int64
	Status         string
}

type TransactionHistoryResult struct {
	Transactions []TransactionInfo
	NextCursor   string
	HasMore      bool
}

type TransactionInfo struct {
	ID        string
	Reference string
	Amount    int64
	Type      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
