package service

import (
	"fmt"
	"github.com/barnettt/banking/db"
	"github.com/barnettt/banking/domain"
	"github.com/barnettt/banking/dto"
	"github.com/barnettt/banking/exceptions"
	"github.com/golang-sql/civil"
	"time"
)

type TransactionService interface {
	NewTransaction(*dto.TransactionRequest) (*dto.TransactionResponse, *exceptions.AppError)
}

type DefaultTransactionService struct {
	repository         domain.TransactionRepository
	accountRepository  domain.AccountRepository
	transactionManager db.TxManager
}

func (defaultTransactionService DefaultTransactionService) NewTransaction(transactionRequest *dto.TransactionRequest) (*dto.TransactionResponse, *exceptions.AppError) {
	transaction := getDomainTransaction(*transactionRequest)
	account, err := defaultTransactionService.accountRepository.GetAccount(transactionRequest.AccountId)
	if err != nil {
		return nil, exceptions.NewDatabaseError(err.Message)
	}
	err = dto.Validate(account, transactionRequest)
	if err != nil {
		return nil, err
	}
	if transaction.TransactionType == "withdrawal" {
		account.Amount = account.Amount - transaction.Amount
	} else if transaction.TransactionType == "deposit" {
		account.Amount = account.Amount + transaction.Amount
	} else {
		return nil, exceptions.NewValidationError("Unknown account type")
	}
	tx, err := defaultTransactionService.transactionManager.StartTransaction()
	err = doAccountUpdate(defaultTransactionService, account)
	if err != nil {
		return nil, exceptions.NewDatabaseError("Unknown account type")
	}

	transactionResponse, err := defaultTransactionService.repository.NewTransaction(transaction)
	if err != nil {
		tx.Rollback()
		return nil, exceptions.NewDatabaseError(err.Message)
	}
	// could find the account by id and update the transaction response with
	// value from db.
	transactionResponse.Balance = account.Amount
	transactionResponse.TransactionAmount = transactionRequest.Amount
	transactionResponse.TransactionType = transactionRequest.TransactionType
	transactionResponse.TransactionDate = transaction.TransactionDate
	// update the account with the transaction
	return transactionResponse, nil
}

func doAccountUpdate(defaultTransactionService DefaultTransactionService, account *dto.AccountRequest) *exceptions.AppError {

	err := defaultTransactionService.accountRepository.UpdateAccount(*getAccount(account))
	if err != nil {
		return exceptions.NewDatabaseError("Unable To update account")
	}
	return nil

}

func getDomainTransaction(request dto.TransactionRequest) domain.Transaction {
	dateTime := civil.DateTimeOf(time.Now())
	dt := domain.Transaction{
		AccountId:       request.AccountId,
		Amount:          request.Amount,
		TransactionType: request.TransactionType,
		TransactionDate: dateTime.String(),
	}
	return dt
}
func getAccount(request *dto.AccountRequest) *domain.Account {
	amount := fmt.Sprintf("%f", request.Amount)
	request.Validate()
	// using time.Now().Format("2021-10-02T11:26:20")
	// kept failing as the date produced seemed to be mangled
	// using civil datetime seems to be ok with mysql db
	dateTime := civil.DateTimeOf(time.Now())
	domAcc := domain.Account{CustomerId: request.CustomerId,
		AccountId:   request.AccountId,
		Amount:      amount,
		OpeningDate: dateTime.String(),
		AccountType: request.AccountType,
		Status:      "1",
	}
	return &domAcc
}

func NewTransactionService(repo domain.TransactionRepositoryDb, accountRepo domain.AccountRepositoryDb, manager db.TxManager) DefaultTransactionService {
	return DefaultTransactionService{repository: repo, accountRepository: accountRepo, transactionManager: manager}
}
