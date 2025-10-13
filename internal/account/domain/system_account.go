package domain

import (
	"time"

	"github.com/google/uuid"
)

type SystemAccount struct {
	ID        string
	LedgerID  string
	Currency  Currency
	Amount    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewSystemAccount(ledgerID string, currency Currency, amount int64) *SystemAccount {
	now := time.Now()
	return &SystemAccount{
		ID:        uuid.New().String(),
		LedgerID:  ledgerID,
		Currency:  currency,
		Amount:    amount,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
