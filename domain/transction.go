package domain

import (
	"github.com/barnettt/banking/dto"
	"github.com/barnettt/banking/exceptions"
)

type Transaction struct {
	TransactionId   string  `db:"transaction_id"`
	AccountId       string  `db:"account_id"`
	Amount          float64 `db:"amount"`
	TransactionType string  `db:"transaction_type"`
	TransactionDate string  `db:"transaction_date"`
}

//go:generate mockgen -destination=../mock/domain/mockTransactionRepository.go -package=domain github.com/barnettt/banking/domain TransactionRepository
type TransactionRepository interface {
	NewTransaction(transaction Transaction) (*dto.TransactionResponse, *exceptions.AppError)
}
