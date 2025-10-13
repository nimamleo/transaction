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
	if err == nil && cachedBalance != nil {
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

	if account.Balance != ledgerBalance {
		if err := s.accountRepo.UpdateBalance(ctx, accountID, ledgerBalance); err != nil {
			return nil, err
		}
		account.Balance = ledgerBalance
		account.UpdatedAt = time.Now()
	}

	if err := s.cache.SetBalance(ctx, accountID, account.Balance, account.UpdatedAt); err != nil {
	}

	return &BalanceInfo{
		Balance:   account.Balance,
		UpdatedAt: account.UpdatedAt,
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

	exists, err := s.accountRepo.TransactionExistsByReference(ctx, reference)
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

	transaction := domain.NewTransaction(accountID, reference, amount, domain.TransactionTypeDeposit)

	if err := s.accountRepo.CreateTransaction(ctx, transaction); err != nil {
		return nil, err
	}

	transferID, err := s.ledger.CreateTransfer(ctx, systemAccount.LedgerID, account.LedgerID, amount)
	if err != nil {
		transaction.Fail()
		s.accountRepo.CreateTransaction(ctx, transaction)
		return nil, err
	}

	newBalance := account.Balance + amount
	if err := s.accountRepo.UpdateBalance(ctx, accountID, newBalance); err != nil {
		transaction.Fail()
		s.accountRepo.CreateTransaction(ctx, transaction)
		return nil, err
	}

	transaction.Complete()
	if err := s.accountRepo.CreateTransaction(ctx, transaction); err != nil {
	}

	if err := s.cache.SetBalance(ctx, accountID, newBalance, time.Now()); err != nil {
	}

	return &DepositResult{
		TransactionID: transaction.ID,
		TransferID:    transferID,
		Amount:        amount,
		NewBalance:    newBalance,
		Status:        string(transaction.Status),
	}, nil
}
