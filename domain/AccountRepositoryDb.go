package domain

import (
	"github.com/barnettt/banking-lib/exceptions"
	"github.com/barnettt/banking-lib/logger"
	"github.com/barnettt/banking/dto"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type AccountRepositoryDb struct {
	client *sqlx.DB
}

func (repository AccountRepositoryDb) Save(inAccount Account) (*dto.NewAccountResponse, *exceptions.AppError) {
	accountInsert := "INSERT INTO ACCOUNTS (customer_id, opening_date, account_type, amount, status) " +
		"VALUES (?, ?, ?, ?, ?)"
	status := inAccount.getDbStatus()

	outAccount, err := repository.client.Exec(accountInsert, inAccount.CustomerId, inAccount.OpeningDate,
		inAccount.AccountType, inAccount.Amount, status)

	if err != nil {
		logger.Error("While inserting new account record : " + err.Error())
		return nil, exceptions.NewDatabaseError("Unexpected database error ")
	}
	var id int64
	id, err = outAccount.LastInsertId()
	if err != nil {
		logger.Error("While retrieving the id for new account : " + err.Error())
		return nil, exceptions.NewDatabaseError("Unexpected database error ")
	}

	res := dto.GetAccountResponse(strconv.Itoa(int(id)))
	return &res, nil
}
func (repository AccountRepositoryDb) GetAccount(accountId string) (*dto.AccountRequest, *exceptions.AppError) {
	accountRetrieve := "SELECT account_id, customer_id, opening_date, account_type, amount, status " +
		"FROM accounts where account_id=?"
	accountResponse := make([]dto.AccountRequest, 0)
	err := repository.client.Select(&accountResponse, accountRetrieve, accountId)
	if err != nil {
		logger.Error(err.Error())
		return nil, exceptions.NewDatabaseError("Unexpected database error")
	}
	return &accountResponse[0], nil
}

// UpdateAccount /* not ideal bit of a hack need the transaction back in the caller
func (repository AccountRepositoryDb) UpdateAccount(account Account) *exceptions.AppError {
	logger.Info("account id : " + account.AccountId + " amount : " + account.Amount)
	updateQuery := "UPDATE accounts SET amount = ? WHERE account_id = ? "
	res, err := repository.client.Exec(updateQuery, account.Amount, account.AccountId)
	if err != nil {
		logger.Error(err.Error())
		return exceptions.NewDatabaseError("Unable to update account")

	}
	id, err := res.LastInsertId()
	if err != nil {
		logger.Error(err.Error())
		return exceptions.NewDatabaseError("Unexpected database error")
	}
	logger.Info("Updated : " + strconv.Itoa(int(id)))
	return nil // don't return anything no need
}

func NewAccountRepositoryDb(client *sqlx.DB) AccountRepositoryDb {
	return AccountRepositoryDb{client}
}
