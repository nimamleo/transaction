package infrastructure

import (
	"context"
	"database/sql"
	"transaction/internal/user/domain"
	"transaction/pkg/genericcode"
	"transaction/pkg/richerror"
)

type apiKeyRepository struct {
	db *sql.DB
}

func NewAPIKeyRepository(db *sql.DB) domain.APIKeyRepository {
	return &apiKeyRepository{db: db}
}

func (r *apiKeyRepository) Create(ctx context.Context, apiKey *domain.APIKey) error {
	query := `
		INSERT INTO api_keys (id, user_id, api_key, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		apiKey.ID,
		apiKey.UserID,
		apiKey.APIKeyHash,
		apiKey.CreatedAt,
		apiKey.ExpiresAt,
	)

	if err != nil {
		return richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to create api key")
	}

	return nil
}

func (r *apiKeyRepository) GetByAPIKey(ctx context.Context, hashedKey string) (*domain.APIKey, error) {
	query := `
		SELECT id, user_id, api_key, created_at, expires_at
		FROM api_keys
		WHERE api_key = $1
	`

	var key domain.APIKey
	err := r.db.QueryRowContext(ctx, query, hashedKey).Scan(
		&key.ID,
		&key.UserID,
		&key.APIKeyHash,
		&key.CreatedAt,
		&key.ExpiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, richerror.NewWithCode(genericcode.Unauthorized, "invalid api key")
	}

	if err != nil {
		return nil, richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to get api key")
	}

	return &key, nil
}

func (r *apiKeyRepository) GetUserIDByAPIKey(ctx context.Context, hashedKey string) (string, error) {

	query := `SELECT user_id FROM api_keys WHERE api_key = $1`

	var userID string
	err := r.db.QueryRowContext(ctx, query, hashedKey).Scan(&userID)

	if err == sql.ErrNoRows {
		return "", richerror.NewWithCode(genericcode.Unauthorized, "invalid api key")
	}

	if err != nil {
		return "", richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to get user id from api key")
	}

	return userID, nil
}
