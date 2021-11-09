package service

import (
	"github.com/barnettt/banking/domain"
	"github.com/barnettt/banking/dto"
	"github.com/barnettt/banking/exceptions"
	"github.com/barnettt/banking/mock/db"
	domain2 "github.com/barnettt/banking/mock/domain"
	"github.com/barnettt/banking/service"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
repository         domain.TransactionRepository
accountRepository  domain.AccountRepository
transactionManager db.TxManager
*/
var transactionService service.TransactionService
var tx *sqlx.Tx
var transactionRepository *domain2.MockTransactionRepository
var transactionManager *db.MockTransactionManager

func setUpForTransactions(t *testing.T) func() {
	controller = gomock.NewController(t)
	transactionRepository = domain2.NewMockTransactionRepository(controller)
	accountRepository = domain2.NewMockAccountRepository(controller)
	transactionManager = db.NewMockTransactionManager(controller)
	transactionService = service.NewTransactionService(transactionRepository, accountRepository, transactionManager)
	defer controller.Finish()
	return func() {
		defer controller.Finish()
	}
}
func Test_should_create_new_transaction_on_account(t *testing.T) {

	// Arrange
	tearDown := setUpForTransactions(t)
	defer tearDown()

	transactionRequest := getATransactionRequest()
	domainTransaction := getADomainTransaction()
	transactionResponse := getATransactionResponse()

	account := getAnAccountRequest()
	domainAccount := getADomainAccount()
	expectations(account, domainAccount, domainTransaction, transactionResponse)
	// Act
	response, _ := transactionService.NewTransaction(transactionRequest)
	// Assert
	assert.NotNil(t, response, "Failed to create new transaction, no response returned")
	assert.Equal(t, transactionResponse, response, "Failed transaction response not matching")
}

func expectations(account *dto.AccountRequest, domainAccount domain.Account, domainTransaction domain.Transaction, transactionResponse *dto.TransactionResponse) {
	accountRepository.EXPECT().GetAccount("1234").Return(account, nil)
	transactionManager.EXPECT().StartTransaction().Return(tx, nil)
	accountRepository.EXPECT().UpdateAccount(domainAccount).Return(nil)
	transactionRepository.EXPECT().NewTransaction(domainTransaction).Return(transactionResponse, nil)
}

func Test_should_throw_unexpected_database_error_when_new_transaction_on_account_called(t *testing.T) {

	// Arrange
	tearDown := setUpForTransactions(t)
	defer tearDown()

	transactionRequest := getATransactionRequest()
	accountRepository.EXPECT().GetAccount("1234").Return(nil, exceptions.NewDatabaseError("Unexpected database error"))
	// Act
	_, err := transactionService.NewTransaction(transactionRequest)
	// Assert
	assert.NotNil(t, err, "Failed to throw exception in error case")
	assert.Equal(t, "Unexpected database error", err.Message, "Failed transaction response not matching")
}

func Test_should_throw_error_when_account_update_fails_when_new_transaction_on_account(t *testing.T) {

	// Arrange
	tearDown := setUpForTransactions(t)
	defer tearDown()
	domainTransaction := getADomainTransaction()
	transactionRequest := getATransactionRequest()
	account := getAnAccountRequest()
	domainAccount := getADomainAccount()

	accountRepository.EXPECT().GetAccount("1234").Return(account, nil)
	transactionManager.EXPECT().StartTransaction().Return(tx, nil)
	accountRepository.EXPECT().UpdateAccount(domainAccount).Return(exceptions.NewDatabaseError("Unable to update account"))
	expectations(account, domainAccount, domainTransaction, nil)
	// Act
	_, err := transactionService.NewTransaction(transactionRequest)
	// Assert
	assert.NotNil(t, err, "Failed to throw error when create new transaction, update account call")
	assert.Equal(t, "Unable to update account", err.Message, "Failed transaction response not matching")
}

func Test_should_throw_error_when_new_transaction_on_account_saved(t *testing.T) {

	// Arrange
	tearDown := setUpForTransactions(t)
	defer tearDown()

	transactionRequest := getATransactionRequest()
	domainTransaction := getADomainTransaction()
	account := getAnAccountRequest()
	domainAccount := getADomainAccount()

	accountRepository.EXPECT().GetAccount("1234").Return(account, nil)
	transactionManager.EXPECT().StartTransaction().Return(tx, nil)
	accountRepository.EXPECT().UpdateAccount(domainAccount).Return(nil)
	transactionManager.EXPECT().StartTransaction().Return(tx, nil)
	transactionManager.EXPECT().RollbackTransaction(tx).Return(exceptions.NewDatabaseError("Unable to save transaction"))
	transactionRepository.EXPECT().NewTransaction(domainTransaction).Return(nil, exceptions.NewDatabaseError("Unable to save transaction"))
	// Act
	_, err := transactionService.NewTransaction(transactionRequest)
	// Assert
	assert.NotNil(t, err, "Failed to create new transaction, no response returned")
	assert.Equal(t, "Unable to save transaction", err.Message, "Failed transaction response not matching")
}

// helper functions for data
func getATransactionRequest() *dto.TransactionRequest {
	return &dto.TransactionRequest{AccountId: "1234", Amount: 500.00, TransactionType: "withdrawal", DateTime: "2021-11-08T22:47:09.440582000"}
}

func getADomainTransaction() domain.Transaction {
	return domain.Transaction{AccountId: "1234", Amount: 500, TransactionType: "withdrawal", TransactionDate: "2021-11-08T22:47:09.440582000"}
}

func getATransactionResponse() *dto.TransactionResponse {
	return &dto.TransactionResponse{
		TransactionId:     "5678",
		TransactionAmount: 500.00,
		TransactionType:   "withdrawal",
		TransactionDate:   "2021-11-08T22:36:55.835011000",
		Balance:           5500.00,
	}
}

func getAnAccountRequest() *dto.AccountRequest {
	return &dto.AccountRequest{
		AccountId:   "1234",
		CustomerId:  "4563",
		OpeningDate: "2021-11-08T22:36:55.835011000",
		AccountType: "withdrawal",
		Amount:      6000.000000,
		Status:      "1",
	}
}

func getADomainAccount() domain.Account {
	return domain.Account{
		AccountId:   "1234",
		CustomerId:  "4563",
		OpeningDate: "2021-11-08T22:36:55.835011000",
		AccountType: "withdrawal",
		Amount:      "5500.000000",
		Status:      "1",
	}
}
