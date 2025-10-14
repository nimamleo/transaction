package application

import (
	"context"
	"time"

	"transaction/internal/account/domain"
)

type MockLock struct{}

var _ domain.Lock = (*MockLock)(nil)

func NewMockLock() *MockLock {
	return &MockLock{}
}

func (m *MockLock) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return true, nil
}

func (m *MockLock) Release(ctx context.Context, key string) error {
	return nil
}

func (m *MockLock) Extend(ctx context.Context, key string, ttl time.Duration) error {
	return nil
}
