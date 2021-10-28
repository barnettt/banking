package domain

import (
	"github.com/barnettt/banking/dto"
	"github.com/barnettt/banking/exceptions"
)

type Account struct {
	AccountId   string `db:"account_id"`
	CustomerId  string `db:"customer_id"`
	OpeningDate string `db:"opening_date"`
	AccountType string `db:"account_type"`
	Amount      string
	Status      string
}

func (account Account) getStatus() string {
	status := "active"
	if account.Status == "0" {
		status = "inactive"
	}
	return status
}

func (account Account) getDbStatus() string {
	status := "1"
	if account.Status == "inactive" {
		status = "0"
	}
	return status
}

type AccountRepository interface {
	Save(Account) (*dto.NewAccountResponse, *exceptions.AppError)
	GetAccount(accountId string) (*dto.AccountRequest, *exceptions.AppError)
	UpdateAccount(Account Account) *exceptions.AppError
}
