package account

import "time"

type Response struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	LedgerID string `json:"ledger_id"`
	Currency string `json:"currency"`
	Balance  int64  `json:"balance"`
}

type BalanceResponse struct {
	Balance   int64     `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DepositResponse struct {
	TransactionID string `json:"transaction_id"`
	TransferID    string `json:"transfer_id"`
	Amount        int64  `json:"amount"`
	NewBalance    int64  `json:"new_balance"`
	Status        string `json:"status"`
}
