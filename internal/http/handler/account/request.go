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

type DepositRequest struct {
	Amount    int64  `json:"amount"`
	Reference string `json:"reference"`
}

func (r DepositRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Amount, validation.Required, validation.Min(1)),
		validation.Field(&r.Reference, validation.Required, validation.Length(1, 255)),
	)
}
