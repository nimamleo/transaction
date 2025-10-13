package application

import (
	"context"
	"time"

	"transaction/internal/account/domain"
)

type Service struct {
	accountRepo domain.AccountRepository
	ledger      domain.Ledger
	cache       domain.AccountCache
}

func NewService(accountRepo domain.AccountRepository, ledger domain.Ledger, cache domain.AccountCache) *Service {
	return &Service{
		accountRepo: accountRepo,
		ledger:      ledger,
		cache:       cache,
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
