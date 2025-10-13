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

	user := httpcontext.GetUser(c)
	if user == nil {
		return stdresponse.SendHttpResponse(c, "user not authenticated")
	}

	createdAccount, err := h.accountService.CreateAccount(c.Request().Context(), user.ID, req.Currency)
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
