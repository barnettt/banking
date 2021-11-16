package domain

import (
	"github.com/barnettt/banking-lib/exceptions"
	"github.com/barnettt/banking-lib/logger"
	"github.com/barnettt/banking/dto"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type TransactionRepositoryDb struct {
	client *sqlx.DB
}

func (repository TransactionRepositoryDb) NewTransaction(transaction Transaction) (*dto.TransactionResponse, *exceptions.AppError) {
	transactionInsert := "INSERT INTO TRANSACTIONS (account_id, amount, transaction_type, transaction_date) values (?, ?, ?, ?)"
	outResponse, err := repository.client.Exec(transactionInsert, transaction.AccountId,
		transaction.Amount, transaction.TransactionType, transaction.TransactionDate,
	)
	if err != nil {
		logger.Error(err.Error())
		return nil, exceptions.NewDatabaseError("Unexpected database error")
	}
	id, err := outResponse.LastInsertId()
	if err != nil {
		logger.Error(err.Error())
		return nil, exceptions.NewDatabaseError("Unexpected database error")
	}
	transactionResponse := dto.GetTransactionResponse(strconv.Itoa(int(id)))
	return &transactionResponse, nil
}

func NewTransactionRepositoryDb(dbClient *sqlx.DB) TransactionRepositoryDb {
	return TransactionRepositoryDb{dbClient}
}
