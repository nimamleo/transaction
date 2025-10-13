package infrastructure

import (
	"context"
	"database/sql"

	"transaction/internal/account/domain"
	"transaction/pkg/genericcode"
	"transaction/pkg/richerror"
)

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) domain.AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *domain.Account) error {
	query := `
		INSERT INTO accounts (id, user_id, ledger_id, currency, balance, version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		account.ID,
		account.UserID,
		account.LedgerID,
		account.Currency.String(),
		account.Balance,
		account.Version,
		account.CreatedAt,
		account.UpdatedAt,
	)

	if err != nil {
		return richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to create account")
	}

	return nil
}

func (r *accountRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM accounts WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to delete account")
	}

	return nil
}

func (r *accountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	query := `
		SELECT id, user_id, ledger_id, currency, balance, version, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`

	var account domain.Account
	var currencyStr string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.LedgerID,
		&currencyStr,
		&account.Balance,
		&account.Version,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, richerror.WrapWithCode(err, genericcode.NotFound, "account not found")
		}
		return nil, richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to fetch account")
	}

	account.Currency = domain.Currency(currencyStr)
	return &account, nil
}

func (r *accountRepository) UpdateBalance(ctx context.Context, id string, balance int64) error {
	query := `
		UPDATE accounts
		SET balance = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, balance, id)
	if err != nil {
		return richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to update balance")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return richerror.NewWithCode(genericcode.NotFound, "account not found")
	}

	return nil
}

func (r *accountRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Account, error) {
	query := `
		SELECT id, user_id, ledger_id, currency, balance, version, created_at, updated_at
		FROM accounts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to fetch accounts")
	}
	defer rows.Close()

	var accounts []*domain.Account
	for rows.Next() {
		var account domain.Account
		var currencyStr string

		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.LedgerID,
			&currencyStr,
			&account.Balance,
			&account.Version,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to scan account")
		}

		account.Currency = domain.Currency(currencyStr)
		accounts = append(accounts, &account)
	}

	if err := rows.Err(); err != nil {
		return nil, richerror.WrapWithCode(err, genericcode.InternalServerError, "error iterating accounts")
	}

	return accounts, nil
}
