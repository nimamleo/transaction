package domain

import (
	"transaction/pkg/genericcode"
	"transaction/pkg/richerror"
)

var (
	ErrUserNotFound       = richerror.NewWithCode(genericcode.NotFound, "user not found")
	ErrEmailAlreadyExists = richerror.NewWithCode(genericcode.Conflict, "email already exists")
)
