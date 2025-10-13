package account

import (
	"transaction/internal/account/application"
	"transaction/pkg/genericcode"
	"transaction/pkg/httpcontext"
	"transaction/pkg/stdresponse"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	accountService *application.Service
}

func NewHandler(accountService *application.Service) *Handler {
	return &Handler{
		accountService: accountService,
	}
}

func (h *Handler) CreateAccount(c echo.Context) error {
	var req CreateAccountRequest

	if err := c.Bind(&req); err != nil {
		return stdresponse.SendHttpResponse(c, err)
	}

	if err := req.Validate(); err != nil {
		return stdresponse.SendHttpResponse(c, err.Error())
	}

	createdAccount, err := h.accountService.CreateAccount(c.Request().Context(), req.UserID, req.Currency)
	if err != nil {
		return stdresponse.SendHttpResponse(c, err)
	}

	return stdresponse.SendHttpResponse(c, genericcode.OK, ToResponse(createdAccount))
}

func (h *Handler) GetAccounts(c echo.Context) error {
	user := httpcontext.GetUser(c)
	if user == nil {
		return stdresponse.SendHttpResponse(c, "user not authenticated")
	}

	accounts, err := h.accountService.GetUserAccounts(c.Request().Context(), user.ID)
	if err != nil {
		return stdresponse.SendHttpResponse(c, err)
	}

	return stdresponse.SendHttpResponse(c, genericcode.OK, ToResponseList(accounts))
}

func (h *Handler) GetAccountBalance(c echo.Context) error {
	accountID := c.Param("id")

	balanceInfo, err := h.accountService.GetAccountBalance(c.Request().Context(), accountID)
	if err != nil {
		return stdresponse.SendHttpResponse(c, err)
	}

	response := BalanceResponse{
		Balance:   balanceInfo.Balance,
		UpdatedAt: balanceInfo.UpdatedAt,
	}

	return stdresponse.SendHttpResponse(c, genericcode.OK, response)
}

func (h *Handler) Deposit(c echo.Context) error {
	accountID := c.Param("id")

	var req DepositRequest
	if err := c.Bind(&req); err != nil {
		return stdresponse.SendHttpResponse(c, err)
	}

	if err := req.Validate(); err != nil {
		return stdresponse.SendHttpResponse(c, err.Error())
	}

	result, err := h.accountService.Deposit(c.Request().Context(), accountID, req.Reference, req.Amount)
	if err != nil {
		return stdresponse.SendHttpResponse(c, err)
	}

	response := DepositResponse{
		TransactionID: result.TransactionID,
		TransferID:    result.TransferID,
		Amount:        result.Amount,
		NewBalance:    result.NewBalance,
		Status:        result.Status,
	}

	return stdresponse.SendHttpResponse(c, genericcode.OK, response)
}

func (h *Handler) Transfer(c echo.Context) error {
	var req TransferRequest
	if err := c.Bind(&req); err != nil {
		return stdresponse.SendHttpResponse(c, err)
	}

	if err := req.Validate(); err != nil {
		return stdresponse.SendHttpResponse(c, err.Error())
	}

	result, err := h.accountService.Transfer(c.Request().Context(), req.FromAccountID, req.ToAccountID, req.Reference, req.Amount)
	if err != nil {
		return stdresponse.SendHttpResponse(c, err)
	}

	response := TransferResponse{
		TransferID:     result.TransferID,
		FromAccountID:  result.FromAccountID,
		ToAccountID:    result.ToAccountID,
		Amount:         result.Amount,
		FromNewBalance: result.FromNewBalance,
		ToNewBalance:   result.ToNewBalance,
		Status:         result.Status,
	}

	return stdresponse.SendHttpResponse(c, genericcode.OK, response)
}

func (h *Handler) GetAccountTransactionHistory(c echo.Context) error {
	accountID := c.Param("id")

	var req TransactionHistoryRequest
	if err := c.Bind(&req); err != nil {
		return stdresponse.SendHttpResponse(c, err)
	}

	if err := req.Validate(); err != nil {
		return stdresponse.SendHttpResponse(c, err.Error())
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	result, err := h.accountService.GetAccountTransactionHistory(c.Request().Context(), accountID, req.Limit, req.After)
	if err != nil {
		return stdresponse.SendHttpResponse(c, err)
	}

	transactions := make([]TransactionResponse, len(result.Transactions))
	for i, tx := range result.Transactions {
		transactions[i] = TransactionResponse{
			ID:        tx.ID,
			Reference: tx.Reference,
			Amount:    tx.Amount,
			Type:      tx.Type,
			Status:    tx.Status,
			CreatedAt: tx.CreatedAt,
			UpdatedAt: tx.UpdatedAt,
		}
	}

	response := TransactionHistoryResponse{
		Transactions: transactions,
		NextCursor:   result.NextCursor,
		HasMore:      result.HasMore,
	}

	return stdresponse.SendHttpResponse(c, genericcode.OK, response)
}
