package dto

import (
	"github.com/barnettt/banking/exceptions"
)

type AccountRequest struct {
	AccountId   string  `json:"account_id"  xml:"account_id" db:"account_id"`
	CustomerId  string  `json:"customer_d"  xml:"customer_id" db:"customer_id"`
	OpeningDate string  `json:"opening_date" xml:"opening_date" db:"opening_date"`
	AccountType string  `json:"account_type" xml:"account_type" db:"account_type"`
	Amount      float64 `json:"amount" xml:"amount"`
	Status      string  `json:"status" xml:"status"`
}

func (account AccountRequest) Validate() *exceptions.AppError {
	if account.Amount < 5000 {
		return exceptions.NewValidationError("Amount is below the required minimum of 5000 ")
	}
	if account.AccountType != "checking" || account.AccountType != "saving" {
		return exceptions.NewValidationError("Invalid Account type,  must be one of checking or saving ")
	}
	return nil
}
