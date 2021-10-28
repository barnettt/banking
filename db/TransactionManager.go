package db

import (
	"github.com/barnettt/banking/exceptions"
	"github.com/jmoiron/sqlx"
)

type TxManager struct {
	DB *sqlx.DB
}

type TransactionManager interface {
	StartTransaction() (*sqlx.Tx, *exceptions.AppError)
	RollbackTransaction(tx *sqlx.Tx) *exceptions.AppError
}

func (transactionManager TxManager) StartTransaction() (*sqlx.Tx, *exceptions.AppError) {
	tx, err := transactionManager.DB.Beginx()
	if err != nil {
		return nil, exceptions.NewDatabaseError("Unable to start transaction")
	}
	return tx, nil
}

func NewTxManager(client *sqlx.DB) TxManager {
	return TxManager{DB: client}
}
