package domain

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        string
	UserID    string
	LedgerID  string
	Currency  Currency
	Balance   int64
	Version   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAccount(userID string, currency Currency) *Account {
	now := time.Now()
	return &Account{
		ID:        uuid.New().String(),
		UserID:    userID,
		Currency:  currency,
		Balance:   0,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
