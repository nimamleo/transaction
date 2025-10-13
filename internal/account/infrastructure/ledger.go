package infrastructure

import (
	"context"
	"fmt"

	"transaction/internal/account/domain"
	"transaction/pkg/genericcode"
	"transaction/pkg/richerror"
	"transaction/pkg/tigerbeetle"

	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"
)

const (
	LedgerID = 1
)

type ledger struct {
	client *tigerbeetle.Client
}

func NewLedger(client *tigerbeetle.Client) domain.Ledger {
	return &ledger{client: client}
}

func (l *ledger) CreateAccount(ctx context.Context, currency domain.Currency) (string, error) {
	tbID := types.ID()

	accounts := []types.Account{
		{
			ID:          tbID,
			UserData128: types.ToUint128(0),
			UserData64:  0,
			UserData32:  0,
			Ledger:      LedgerID,
			Code:        currency.Code(),
			Flags:       0,
			Timestamp:   0,
		},
	}

	results, err := l.client.GetClient().CreateAccounts(accounts)
	if err != nil {
		return "", richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to create ledger account")
	}

	if len(results) > 0 {
		return "", richerror.NewWithCode(genericcode.InternalServerError, fmt.Sprintf("ledger account creation failed: %v", results[0]))
	}

	return uint128ToString(tbID), nil
}

func (l *ledger) GetBalance(ctx context.Context, ledgerID string) (int64, error) {
	id, err := stringToUint128(ledgerID)
	if err != nil {
		return 0, richerror.WrapWithCode(err, genericcode.BadRequest, "invalid ledger ID format")
	}

	accounts, err := l.client.GetClient().LookupAccounts([]types.Uint128{id})
	if err != nil {
		return 0, richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to lookup account in ledger")
	}

	if len(accounts) == 0 {
		return 0, richerror.NewWithCode(genericcode.NotFound, "account not found in ledger")
	}

	account := accounts[0]
	debitsLow := account.DebitsPosted[len(account.DebitsPosted)-1]
	creditsLow := account.CreditsPosted[len(account.CreditsPosted)-1]

	balance := int64(debitsLow) - int64(creditsLow)
	return balance, nil
}

func (l *ledger) CreateTransfer(ctx context.Context, fromLedgerID, toLedgerID string, amount int64) (string, error) {
	fromID, err := stringToUint128(fromLedgerID)
	if err != nil {
		return "", richerror.WrapWithCode(err, genericcode.BadRequest, "invalid from ledger ID format")
	}

	toID, err := stringToUint128(toLedgerID)
	if err != nil {
		return "", richerror.WrapWithCode(err, genericcode.BadRequest, "invalid to ledger ID format")
	}

	transferID := types.ID()

	transfers := []types.Transfer{
		{
			ID:              transferID,
			DebitAccountID:  fromID,
			CreditAccountID: toID,
			Amount:          types.ToUint128(uint64(amount)),
			Ledger:          LedgerID,
			Code:            1,
			Flags:           0,
			Timestamp:       0,
		},
	}

	results, err := l.client.GetClient().CreateTransfers(transfers)
	if err != nil {
		return "", richerror.WrapWithCode(err, genericcode.InternalServerError, "failed to create ledger transfer")
	}

	if len(results) > 0 {
		return "", richerror.NewWithCode(genericcode.InternalServerError, fmt.Sprintf("ledger transfer creation failed: %v", results[0]))
	}

	return uint128ToString(transferID), nil
}
