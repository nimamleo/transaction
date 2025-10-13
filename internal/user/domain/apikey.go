package domain

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"transaction/pkg/hash"

	"github.com/google/uuid"
)

type APIKey struct {
	ID          string
	UserID      string
	APIKeyHash  string
	PlainAPIKey string
	CreatedAt   time.Time
	ExpiresAt   *time.Time
}

func NewAPIKey(userID string) (*APIKey, error) {
	plainKey, err := generateAPIKey()
	if err != nil {
		return nil, err
	}

	return &APIKey{
		ID:          uuid.New().String(),
		UserID:      userID,
		APIKeyHash:  hash.Hash(plainKey),
		PlainAPIKey: plainKey,
		CreatedAt:   time.Now(),
		ExpiresAt:   nil,
	}, nil
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
