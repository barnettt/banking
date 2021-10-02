package domain

import (
	"github.com/barnettt/banking/dto"
	"github.com/barnettt/banking/exceptions"
	"github.com/barnettt/banking/logger"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type AccountRepositoryDb struct {
	client *sqlx.DB
}

func (repository AccountRepositoryDb) Save(inAccount Account) (*dto.NewAccountResponse, *exceptions.AppError) {
	accountInsert := "INSERT INTO ACCOUNTS (customer_id, opening_date, account_type, amount, status) " +
		"VALUES (?, ?, ?, ?, ?)"
		// "VALUES (?, STR_TO_DATE(?,'%d/%m/%Y %H:%i%s'), ?, ?, ?)"
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

func NewAccountRepositoryDb(client *sqlx.DB) AccountRepositoryDb {
	return AccountRepositoryDb{client}
}
