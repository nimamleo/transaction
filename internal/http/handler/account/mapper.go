package account

import "transaction/internal/account/domain"

func ToResponse(account *domain.Account) Response {
	return Response{
		ID:       account.ID,
		UserID:   account.UserID,
		LedgerID: account.LedgerID,
		Currency: account.Currency.String(),
		Balance:  account.Balance,
	}
}

func ToResponseList(accounts []*domain.Account) []Response {
	responses := make([]Response, len(accounts))
	for i, account := range accounts {
		responses[i] = ToResponse(account)
	}
	return responses
}
