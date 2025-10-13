package domain

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type APIKeyRepository interface {
	Create(ctx context.Context, apiKey *APIKey) error
	GetByAPIKey(ctx context.Context, apiKey string) (*APIKey, error)
	GetUserIDByAPIKey(ctx context.Context, apiKey string) (string, error)
}
