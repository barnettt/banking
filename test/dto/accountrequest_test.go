package dto

import (
	"github.com/barnettt/banking/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_should_open_account_with_minimum_funds(t *testing.T) {
	account := dto.AccountRequest{
		AccountId:   "1234",
		AccountType: "checking",
		Amount:      5000,
	}

	err := dto.AccountRequest.Validate(account)

	assert.Nil(t, err, "failed open account with minimum funds test ")
}

func Test_should_not_open_account_when_not_minimum_funds_met(t *testing.T) {
	account := dto.AccountRequest{
		AccountId:   "1234",
		AccountType: "checking",
		Amount:      3000,
	}

	err := dto.AccountRequest.Validate(account)

	assert.Equal(t, err.AsMessage().Message, "Amount is below the required minimum of 5000 ", "failed minimum account funds not met funds test ")
}

func Test_should_not_open_account_when_invalid_account_type(t *testing.T) {
	account := dto.AccountRequest{
		AccountId:   "1234",
		AccountType: "checking333",
		Amount:      5000,
	}

	err := dto.AccountRequest.Validate(account)

	assert.Equal(t, err.AsMessage().Message, "Invalid Account type,  must be one of checking or saving ", "failed minimum account funds not met funds test ")
}
