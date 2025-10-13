package account

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateAccountRequest struct {
	Currency string `json:"currency"`
}

func (r CreateAccountRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Currency, validation.Required, validation.In("USD", "EUR", "GBP")),
	)
}
