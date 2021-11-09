package db

import (
	"github.com/barnettt/banking/exceptions"
	"github.com/jmoiron/sqlx"
)

type TxManager struct {
	DB *sqlx.DB
}

//go:generate mockgen -destination=../mock/db/mockTransactionManager.go -package=db github.com/barnettt/banking/db TransactionManager
type TransactionManager interface {
	StartTransaction() (*sqlx.Tx, *exceptions.AppError)
	RollbackTransaction(tx *sqlx.Tx) *exceptions.AppError
}

func (transactionManager TxManager) StartTransaction() (*sqlx.Tx, *exceptions.AppError) {
	tx, err := transactionManager.DB.Beginx()
	if err != nil {
		return nil, exceptions.NewDatabaseError("Error Unable to start transaction")
	}
	return tx, nil
}

func (transactionManager TxManager) RollbackTransaction(tx *sqlx.Tx) *exceptions.AppError {
	if err := transactionManager.RollbackTransaction(tx); err != nil {
		return exceptions.NewDatabaseError("Error Unable to rollback transaction")
	}
	return nil
}

func NewTxManager(client *sqlx.DB) TxManager {
	return TxManager{DB: client}
}
