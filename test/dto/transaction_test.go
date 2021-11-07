package dto

import (
	"github.com/barnettt/banking/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_should_be_valid_transaction_type(t *testing.T) {
	// Arrange
	ac := dto.AccountRequest{
		Amount:      500,
		AccountType: "withdrawal",
	}
	tr := dto.TransactionRequest{
		Amount:          500,
		TransactionType: "withdrawal",
	}
	// Act
	err := dto.Validate(&ac, &tr)
	// Assert
	assert.Nil(t, err, "Failed valid transaction type test")
}

func Test_should_be_invalid_transaction_type(t *testing.T) {
	// Arrange
	ac := dto.AccountRequest{
		Amount:      500,
		AccountType: "withdrawalsss",
	}
	tr := dto.TransactionRequest{
		Amount:          500,
		TransactionType: "withdrawalsss",
	}
	// Act
	err := dto.Validate(&ac, &tr)
	// Assert
	assert.Equal(t, err.AsMessage().Message, "Transaction type must be one of 'withdrawal' or 'deposit'", "Failed invalid transaction type test")
}

func Test_should_be_non_negative_transaction_amount(t *testing.T) {
	// Arrange
	ac := dto.AccountRequest{
		Amount:      -500,
		AccountType: "withdrawal",
	}
	tr := dto.TransactionRequest{
		Amount:          -500,
		TransactionType: "withdrawal",
	}
	// Act
	err := dto.Validate(&ac, &tr)
	// Assert
	assert.Equal(t, err.AsMessage().Message, "Cannot have a negative transaction amount", "failed, negative transaction amount test")
}

func Test_should_throw_error_for_insufficient_funds(t *testing.T) {
	// Arrange
	ac := dto.AccountRequest{
		Amount:      200,
		AccountType: "withdrawal",
	}
	tr := dto.TransactionRequest{
		Amount:          500,
		TransactionType: "withdrawal",
	}
	// Act
	err := dto.Validate(&ac, &tr)
	// Assert
	assert.Equal(t, err.AsMessage().Message, "Insufficient funds, cannot complete transaction", "failed insufficient funds in account test")
}
