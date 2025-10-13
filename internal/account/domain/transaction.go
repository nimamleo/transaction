package domain

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID        string
	AccountID string
	Reference string
	Amount    int64
	Type      TransactionType
	Status    TransactionStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeTransfer TransactionType = "transfer"
	TransactionTypeWithdraw TransactionType = "withdraw"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

func NewTransaction(accountID, reference string, amount int64, transactionType TransactionType) *Transaction {
	now := time.Now()
	return &Transaction{
		ID:        uuid.New().String(),
		AccountID: accountID,
		Reference: reference,
		Amount:    amount,
		Type:      transactionType,
		Status:    TransactionStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (t *Transaction) Complete() {
	t.Status = TransactionStatusCompleted
	t.UpdatedAt = time.Now()
}

func (t *Transaction) Fail() {
	t.Status = TransactionStatusFailed
	t.UpdatedAt = time.Now()
}
