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

func (r *accountRepository) CreateSystemAccount(ctx context.Context, systemAccount *domain.SystemAccount) error {
	query := `
		INSERT INTO system_accounts (id, ledger_id, currency, amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		systemAccount.ID,
		systemAccount.LedgerID,
		systemAccount.Currency.String(),
		systemAccount.Amount,
		systemAccount.CreatedAt,
		systemAccount.UpdatedAt,
	)

	if err != nil {
		return richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to create system account")
	}

	return nil
}

func (r *accountRepository) GetSystemAccountByCurrency(ctx context.Context, currency domain.Currency) (*domain.SystemAccount, error) {
	query := `
		SELECT id, ledger_id, currency, amount, created_at, updated_at
		FROM system_accounts
		WHERE currency = $1
	`

	var systemAccount domain.SystemAccount
	var currencyStr string

	err := r.db.QueryRowContext(ctx, query, currency.String()).Scan(
		&systemAccount.ID,
		&systemAccount.LedgerID,
		&currencyStr,
		&systemAccount.Amount,
		&systemAccount.CreatedAt,
		&systemAccount.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, richerror.WrapWithCode(err, genericcode.NotFound, "system account not found")
		}
		return nil, richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to fetch system account")
	}

	systemAccount.Currency = domain.Currency(currencyStr)
	return &systemAccount, nil
}

func (r *accountRepository) SystemAccountExistsByCurrency(ctx context.Context, currency domain.Currency) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM system_accounts WHERE currency = $1
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, currency.String()).Scan(&exists)
	if err != nil {
		return false, richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to check system account existence")
	}

	return exists, nil
}

func (r *accountRepository) GetTransactionByReference(ctx context.Context, reference string) (*domain.Transaction, error) {
	query := `
		SELECT id, account_id, reference, amount, type, status, created_at, updated_at
		FROM transactions
		WHERE reference = $1
	`

	var transaction domain.Transaction
	var typeStr, statusStr string

	err := r.db.QueryRowContext(ctx, query, reference).Scan(
		&transaction.ID,
		&transaction.AccountID,
		&transaction.Reference,
		&transaction.Amount,
		&typeStr,
		&statusStr,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, richerror.WrapWithCode(err, genericcode.NotFound, "transaction not found")
		}
		return nil, richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to fetch transaction")
	}

	transaction.Type = domain.TransactionType(typeStr)
	transaction.Status = domain.TransactionStatus(statusStr)
	return &transaction, nil
}

func (r *accountRepository) TransactionExistsByReference(ctx context.Context, reference string, accountID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM transactions WHERE reference = $1 AND account_id = $2
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, reference, accountID).Scan(&exists)
	if err != nil {
		return false, richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to check transaction existence")
	}

	return exists, nil
}

func (r *accountRepository) CreateTransactionAndUpdateBalance(ctx context.Context, transaction *domain.Transaction, accountID string, newBalance int64) (*domain.Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to begin transaction")
	}
	defer tx.Rollback()

	transaction.Complete()

	createTransactionQuery := `
		INSERT INTO transactions (id, account_id, reference, amount, type, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = tx.ExecContext(ctx, createTransactionQuery,
		transaction.ID,
		transaction.AccountID,
		transaction.Reference,
		transaction.Amount,
		string(transaction.Type),
		string(transaction.Status),
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)
	if err != nil {
		return nil, richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to create transaction")
	}

	updateBalanceQuery := `
		UPDATE accounts
		SET balance = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := tx.ExecContext(ctx, updateBalanceQuery, newBalance, accountID)
	if err != nil {
		return nil, richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to update balance")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return nil, richerror.NewWithCode(genericcode.NotFound, "account not found")
	}

	if err := tx.Commit(); err != nil {
		return nil, richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to commit transaction")
	}

	return transaction, nil
}

func (r *accountRepository) CreateTransferTransactions(ctx context.Context, fromAccountID, toAccountID, reference string, amount int64, fromNewBalance, toNewBalance int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to begin transaction")
	}
	defer tx.Rollback()

	fromTransaction := domain.NewTransaction(fromAccountID, reference, amount, domain.TransactionTypeTransfer)
	fromTransaction.Complete()

	toTransaction := domain.NewTransaction(toAccountID, reference, amount, domain.TransactionTypeTransfer)
	toTransaction.Complete()

	createTransactionQuery := `
		INSERT INTO transactions (id, account_id, reference, amount, type, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = tx.ExecContext(ctx, createTransactionQuery,
		fromTransaction.ID,
		fromTransaction.AccountID,
		fromTransaction.Reference,
		-fromTransaction.Amount,
		string(fromTransaction.Type),
		string(fromTransaction.Status),
		fromTransaction.CreatedAt,
		fromTransaction.UpdatedAt,
	)
	if err != nil {
		return richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to create from transaction")
	}

	_, err = tx.ExecContext(ctx, createTransactionQuery,
		toTransaction.ID,
		toTransaction.AccountID,
		toTransaction.Reference,
		toTransaction.Amount,
		string(toTransaction.Type),
		string(toTransaction.Status),
		toTransaction.CreatedAt,
		toTransaction.UpdatedAt,
	)
	if err != nil {
		return richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to create to transaction")
	}

	updateBalanceQuery := `
		UPDATE accounts
		SET balance = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := tx.ExecContext(ctx, updateBalanceQuery, fromNewBalance, fromAccountID)
	if err != nil {
		return richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to update from account balance")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return richerror.NewWithCode(genericcode.NotFound, "from account not found")
	}

	result, err = tx.ExecContext(ctx, updateBalanceQuery, toNewBalance, toAccountID)
	if err != nil {
		return richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to update to account balance")
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return richerror.NewWithCode(genericcode.NotFound, "to account not found")
	}

	if err := tx.Commit(); err != nil {
		return richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to commit transaction")
	}

	return nil
}
