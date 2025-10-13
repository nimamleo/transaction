package domain

import (
	"transaction/pkg/genericcode"
	"transaction/pkg/richerror"
)

var (
	ErrAccountNotFound          = richerror.NewWithCode(genericcode.NotFound, "account not found")
	ErrInsufficientFunds        = richerror.NewWithCode(genericcode.BadRequest, "insufficient funds")
	ErrInvalidCurrency          = richerror.NewWithCode(genericcode.BadRequest, "invalid currency")
	ErrUserNotFound             = richerror.NewWithCode(genericcode.NotFound, "user not found")
	ErrAccountAlreadyExists     = richerror.NewWithCode(genericcode.Conflict, "account already exists")
	ErrInvalidAmount            = richerror.NewWithCode(genericcode.BadRequest, "invalid amount")
	ErrLockAcquisitionFailed    = richerror.NewWithCode(genericcode.InternalServerError, "failed to acquire lock")
	ErrTransactionAlreadyExists = richerror.NewWithCode(genericcode.BadRequest, "transaction with this reference already exists")
)
