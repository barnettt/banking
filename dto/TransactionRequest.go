package dto

import (
	"github.com/barnettt/banking-lib/exceptions"
)

type TransactionRequest struct {
	AccountId       string  `json:"account_id"  xml:"account_id"`
	Amount          float64 `json:"amount"  xml:"amount"`
	TransactionType string  `json:"transaction_type" xml:"transaction_type"`
	DateTime        string  `json:"_"  xml:"_"`
}

func Validate(account *AccountRequest, transaction *TransactionRequest) *exceptions.AppError {

	if transaction.TransactionType != "deposit" && transaction.TransactionType != "withdrawal" {
		return exceptions.NewValidationError("Transaction type must be one of 'withdrawal' or 'deposit'")
	} else if account.Amount < transaction.Amount && transaction.TransactionType == "withdrawal" {
		return exceptions.NewValidationError("Insufficient funds, cannot complete transaction")
	} else if transaction.Amount < 0 {
		return exceptions.NewValidationError("Cannot have a negative transaction amount")
	}
	return nil
}
