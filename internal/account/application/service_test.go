package application

import (
	"context"
	"testing"
	"time"

	"transaction/internal/account/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(ctx context.Context, account *domain.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Account, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) UpdateBalance(ctx context.Context, id string, balance int64) error {
	args := m.Called(ctx, id, balance)
	return args.Error(0)
}

func (m *MockAccountRepository) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockAccountRepository) GetAccountTransactions(ctx context.Context, accountID string, limit int, after string) ([]*domain.Transaction, error) {
	args := m.Called(ctx, accountID, limit, after)
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockAccountRepository) TransactionExistsByReference(ctx context.Context, reference, accountID string) (bool, error) {
	args := m.Called(ctx, reference, accountID)
	return args.Bool(0), args.Error(1)
}

func (m *MockAccountRepository) CreateTransactionAndUpdateBalance(ctx context.Context, transaction *domain.Transaction, accountID string, newBalance int64) (*domain.Transaction, error) {
	args := m.Called(ctx, transaction, accountID, newBalance)
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockAccountRepository) CreateTransferTransactions(ctx context.Context, fromAccountID, toAccountID, reference string, amount, fromNewBalance, toNewBalance int64) error {
	args := m.Called(ctx, fromAccountID, toAccountID, reference, amount, fromNewBalance, toNewBalance)
	return args.Error(0)
}

func (m *MockAccountRepository) SystemAccountExistsByCurrency(ctx context.Context, currency domain.Currency) (bool, error) {
	args := m.Called(ctx, currency)
	return args.Bool(0), args.Error(1)
}

func (m *MockAccountRepository) GetSystemAccountByCurrency(ctx context.Context, currency domain.Currency) (*domain.SystemAccount, error) {
	args := m.Called(ctx, currency)
	return args.Get(0).(*domain.SystemAccount), args.Error(1)
}

func (m *MockAccountRepository) CreateSystemAccount(ctx context.Context, account *domain.SystemAccount) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAccountRepository) GetTransactionByReference(ctx context.Context, reference string) (*domain.Transaction, error) {
	args := m.Called(ctx, reference)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

type MockLedger struct {
	mock.Mock
}

func (m *MockLedger) CreateAccount(ctx context.Context, currency domain.Currency) (string, error) {
	args := m.Called(ctx, currency)
	return args.String(0), args.Error(1)
}

func (m *MockLedger) GetBalance(ctx context.Context, ledgerID string) (int64, error) {
	args := m.Called(ctx, ledgerID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLedger) CreateTransfer(ctx context.Context, fromLedgerID, toLedgerID string, amount int64) (string, error) {
	args := m.Called(ctx, fromLedgerID, toLedgerID, amount)
	return args.String(0), args.Error(1)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) GetBalance(ctx context.Context, accountID string) (*domain.BalanceCache, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.BalanceCache), args.Error(1)
}

func (m *MockCache) SetBalance(ctx context.Context, accountID string, balance int64, updatedAt time.Time) error {
	args := m.Called(ctx, accountID, balance, updatedAt)
	return args.Error(0)
}

func TestService_CreateAccount(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAccountRepository{}
	mockLedger := &MockLedger{}
	mockCache := &MockCache{}
	service := NewService(mockRepo, mockLedger, mockCache, NewMockLock())

	userID := "user-123"
	currency := domain.USD
	ledgerID := "ledger-123"

	mockLedger.On("CreateAccount", ctx, currency).Return(ledgerID, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Account")).Return(nil)

	account, err := service.CreateAccount(ctx, userID, string(currency))

	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, userID, account.UserID)
	assert.Equal(t, currency, account.Currency)
	assert.Equal(t, ledgerID, account.LedgerID)

	mockRepo.AssertExpectations(t)
	mockLedger.AssertExpectations(t)
}

func TestService_GetAccountBalance_CacheHit(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAccountRepository{}
	mockLedger := &MockLedger{}
	mockCache := &MockCache{}
	service := NewService(mockRepo, mockLedger, mockCache, NewMockLock())

	accountID := "account-123"
	cachedBalance := &domain.BalanceCache{
		Balance:   1000,
		UpdatedAt: time.Now(),
	}

	mockCache.On("GetBalance", ctx, accountID).Return(cachedBalance, nil)

	balanceInfo, err := service.GetAccountBalance(ctx, accountID)

	assert.NoError(t, err)
	assert.NotNil(t, balanceInfo)
	assert.Equal(t, int64(1000), balanceInfo.Balance)

	mockCache.AssertExpectations(t)
}

func TestService_GetAccountBalance_CacheMiss(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAccountRepository{}
	mockLedger := &MockLedger{}
	mockCache := &MockCache{}
	service := NewService(mockRepo, mockLedger, mockCache, NewMockLock())

	accountID := "account-123"
	ledgerID := "ledger-123"
	account := &domain.Account{
		ID:       accountID,
		LedgerID: ledgerID,
		Balance:  500,
	}

	mockCache.On("GetBalance", ctx, accountID).Return(nil, nil)
	mockRepo.On("GetByID", ctx, accountID).Return(account, nil)
	mockLedger.On("GetBalance", ctx, ledgerID).Return(int64(1000), nil)
	mockCache.On("SetBalance", ctx, accountID, int64(1000), mock.AnythingOfType("time.Time")).Return(nil)

	balanceInfo, err := service.GetAccountBalance(ctx, accountID)

	assert.NoError(t, err)
	assert.NotNil(t, balanceInfo)
	assert.Equal(t, int64(1000), balanceInfo.Balance)

	mockRepo.AssertExpectations(t)
	mockLedger.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestService_Deposit_InvalidAmount(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAccountRepository{}
	mockLedger := &MockLedger{}
	mockCache := &MockCache{}
	service := NewService(mockRepo, mockLedger, mockCache, NewMockLock())

	result, err := service.Deposit(ctx, "account-123", "ref-123", -100)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrInvalidAmount, err)
}

func TestService_Transfer_InvalidAmount(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAccountRepository{}
	mockLedger := &MockLedger{}
	mockCache := &MockCache{}
	service := NewService(mockRepo, mockLedger, mockCache, NewMockLock())

	result, err := service.Transfer(ctx, "from-123", "to-123", "ref-123", -100)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrInvalidAmount, err)
}

func TestService_Transfer_SameAccount(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAccountRepository{}
	mockLedger := &MockLedger{}
	mockCache := &MockCache{}
	service := NewService(mockRepo, mockLedger, mockCache, NewMockLock())

	result, err := service.Transfer(ctx, "account-123", "account-123", "ref-123", 100)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrSameAccountTransfer, err)
}

func TestService_Transfer_InsufficientFunds(t *testing.T) {
	ctx := context.Background()
	mockRepo := &MockAccountRepository{}
	mockLedger := &MockLedger{}
	mockCache := &MockCache{}
	service := NewService(mockRepo, mockLedger, mockCache, NewMockLock())

	fromAccountID := "from-account-123"
	toAccountID := "to-account-123"
	reference := "transfer-ref-123"
	amount := int64(1500)

	fromAccount := &domain.Account{
		ID:       fromAccountID,
		LedgerID: "from-ledger-123",
		Balance:  1000,
		Currency: domain.USD,
	}

	toAccount := &domain.Account{
		ID:       toAccountID,
		LedgerID: "to-ledger-123",
		Balance:  200,
		Currency: domain.USD,
	}

	mockRepo.On("TransactionExistsByReference", ctx, reference, fromAccountID).Return(false, nil)
	mockRepo.On("GetByID", ctx, fromAccountID).Return(fromAccount, nil)
	mockRepo.On("GetByID", ctx, toAccountID).Return(toAccount, nil)

	result, err := service.Transfer(ctx, fromAccountID, toAccountID, reference, amount)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, domain.ErrInsufficientFunds, err)

	mockRepo.AssertExpectations(t)
}
