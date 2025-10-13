package application

import (
	"context"
	"fmt"
	"time"

	"transaction/internal/account/domain"
	"transaction/internal/account/infrastructure"
)

type Service struct {
	accountRepo domain.AccountRepository
	ledger      domain.Ledger
	cache       domain.AccountCache
	lock        *infrastructure.Lock
}

func NewService(accountRepo domain.AccountRepository, ledger domain.Ledger, cache domain.AccountCache, lock *infrastructure.Lock) *Service {
	return &Service{
		accountRepo: accountRepo,
		ledger:      ledger,
		cache:       cache,
		lock:        lock,
	}
}

func (s *Service) CreateAccount(ctx context.Context, userID, currencyStr string) (*domain.Account, error) {
	currency := domain.Currency(currencyStr)

	if !currency.IsValid() {
		return nil, domain.ErrInvalidCurrency
	}

	ledgerID, err := s.ledger.CreateAccount(ctx, currency)
	if err != nil {
		return nil, err
	}

	account := domain.NewAccount(userID, currency)
	account.LedgerID = ledgerID

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *Service) GetUserAccounts(ctx context.Context, userID string) ([]*domain.Account, error) {
	return s.accountRepo.GetByUserID(ctx, userID)
}

func (s *Service) GetAccountBalance(ctx context.Context, accountID string) (*BalanceInfo, error) {
	cachedBalance, err := s.cache.GetBalance(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if cachedBalance != nil {
		return &BalanceInfo{
			Balance:   cachedBalance.Balance,
			UpdatedAt: cachedBalance.UpdatedAt,
		}, nil
	}

	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	ledgerBalance, err := s.ledger.GetBalance(ctx, account.LedgerID)
	if err != nil {
		return nil, err
	}

	updatedAt := time.Now()
	if err := s.cache.SetBalance(ctx, accountID, ledgerBalance, updatedAt); err != nil {
		return nil, err
	}

	return &BalanceInfo{
		Balance:   ledgerBalance,
		UpdatedAt: updatedAt,
	}, nil
}

func (s *Service) InitializeSystemAccount(ctx context.Context, currency domain.Currency, amount int64) error {
	exists, err := s.accountRepo.SystemAccountExistsByCurrency(ctx, currency)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	ledgerID, err := s.ledger.CreateAccount(ctx, currency)
	if err != nil {
		return err
	}

	systemAccount := domain.NewSystemAccount(ledgerID, currency, amount)

	return s.accountRepo.CreateSystemAccount(ctx, systemAccount)
}

func (s *Service) Deposit(ctx context.Context, accountID, reference string, amount int64) (*DepositResult, error) {
	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	lockKey := fmt.Sprintf("deposit:%s:%s", accountID, reference)
	lockTTL := 30 * time.Second

	acquired, err := s.lock.Acquire(ctx, lockKey, lockTTL)
	if err != nil {
		return nil, err
	}
	if !acquired {
		return nil, domain.ErrLockAcquisitionFailed
	}

	defer func() {
		if releaseErr := s.lock.Release(ctx, lockKey); releaseErr != nil {
		}
	}()

	exists, err := s.accountRepo.TransactionExistsByReference(ctx, reference, accountID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrTransactionAlreadyExists
	}

	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	systemAccount, err := s.accountRepo.GetSystemAccountByCurrency(ctx, account.Currency)
	if err != nil {
		return nil, err
	}

	transferID, err := s.ledger.CreateTransfer(ctx, systemAccount.LedgerID, account.LedgerID, amount)
	if err != nil {
		return nil, err
	}

	transaction := domain.NewTransaction(accountID, reference, amount, domain.TransactionTypeDeposit)
	newBalance := account.Balance + amount

	result, err := s.accountRepo.CreateTransactionAndUpdateBalance(ctx, transaction, accountID, newBalance)
	if err != nil {
		return nil, err
	}

	updatedAt := time.Now()
	if err := s.cache.SetBalance(ctx, accountID, newBalance, updatedAt); err != nil {
		return nil, err
	}

	return &DepositResult{
		TransactionID: transaction.ID,
		TransferID:    transferID,
		Amount:        amount,
		NewBalance:    newBalance,
		Status:        string(result.Status),
	}, nil
}

func (s *Service) Transfer(ctx context.Context, fromAccountID, toAccountID, reference string, amount int64) (*TransferResult, error) {
	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	if fromAccountID == toAccountID {
		return nil, domain.ErrSameAccountTransfer
	}

	lockKey := fmt.Sprintf("transfer:%s:%s", fromAccountID, reference)
	lockTTL := 30 * time.Second

	acquired, err := s.lock.Acquire(ctx, lockKey, lockTTL)
	if err != nil {
		return nil, err
	}
	if !acquired {
		return nil, domain.ErrLockAcquisitionFailed
	}

	defer func() {
		if releaseErr := s.lock.Release(ctx, lockKey); releaseErr != nil {
		}
	}()

	exists, err := s.accountRepo.TransactionExistsByReference(ctx, reference, fromAccountID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrTransactionAlreadyExists
	}

	fromAccount, err := s.accountRepo.GetByID(ctx, fromAccountID)
	if err != nil {
		return nil, err
	}

	toAccount, err := s.accountRepo.GetByID(ctx, toAccountID)
	if err != nil {
		return nil, err
	}

	if fromAccount.Currency != toAccount.Currency {
		return nil, domain.ErrCurrencyMismatch
	}

	if fromAccount.Balance < amount {
		return nil, domain.ErrInsufficientFunds
	}

	transferID, err := s.ledger.CreateTransfer(ctx, fromAccount.LedgerID, toAccount.LedgerID, amount)
	if err != nil {
		return nil, err
	}

	fromNewBalance := fromAccount.Balance - amount
	toNewBalance := toAccount.Balance + amount

	err = s.accountRepo.CreateTransferTransactions(ctx, fromAccountID, toAccountID, reference, amount, fromNewBalance, toNewBalance)
	if err != nil {
		return nil, err
	}

	updatedAt := time.Now()
	if err := s.cache.SetBalance(ctx, fromAccountID, fromNewBalance, updatedAt); err != nil {
		return nil, err
	}
	if err := s.cache.SetBalance(ctx, toAccountID, toNewBalance, updatedAt); err != nil {
		return nil, err
	}

	return &TransferResult{
		TransferID:     transferID,
		FromAccountID:  fromAccountID,
		ToAccountID:    toAccountID,
		Amount:         amount,
		FromNewBalance: fromNewBalance,
		ToNewBalance:   toNewBalance,
		Status:         string(domain.TransactionStatusCompleted),
	}, nil
}
