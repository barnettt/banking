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

//go:generate mockgen -destination=../mock/service/mockTransactionService.go -package=service github.com/barnettt/banking/service TransactionService
type TransactionService interface {
	NewTransaction(*dto.TransactionRequest) (*dto.TransactionResponse, *exceptions.AppError)
}

type DefaultTransactionService struct {
	Repository         domain.TransactionRepository
	AccountRepository  domain.AccountRepository
	TransactionManager db.TransactionManager //db.TxManager
}

func (defaultTransactionService DefaultTransactionService) NewTransaction(transactionRequest *dto.TransactionRequest) (*dto.TransactionResponse, *exceptions.AppError) {
	transaction := getDomainTransaction(*transactionRequest)
	account, err := defaultTransactionService.AccountRepository.GetAccount(transactionRequest.AccountId)
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
	tx, err := defaultTransactionService.TransactionManager.StartTransaction()
	err = doAccountUpdate(defaultTransactionService, account)
	if err != nil {
		return nil, exceptions.NewDatabaseError("Unable to update account")
	}

	transactionResponse, err := defaultTransactionService.Repository.NewTransaction(transaction)
	if err != nil {
		ex := defaultTransactionService.TransactionManager.RollbackTransaction(tx)
		if ex != nil {
			return nil, ex
		}
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

	err := defaultTransactionService.AccountRepository.UpdateAccount(*getAccount(account))
	if err != nil {
		return exceptions.NewDatabaseError("Unable To update account")
	}
	return nil

}

func getDomainTransaction(request dto.TransactionRequest) domain.Transaction {
	dateTime := civil.DateTimeOf(time.Now())
	var dateStr = dateTime.String()
	if request.DateTime != "" {
		dateStr = request.DateTime
	}
	dt := domain.Transaction{
		AccountId:       request.AccountId,
		Amount:          request.Amount,
		TransactionType: request.TransactionType,
		TransactionDate: dateStr,
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
	var dateStr = dateTime.String()
	if request.OpeningDate != "" {
		dateStr = request.OpeningDate
	}
	domAcc := domain.Account{CustomerId: request.CustomerId,
		AccountId:   request.AccountId,
		Amount:      amount,
		OpeningDate: dateStr,
		AccountType: request.AccountType,
		Status:      "1",
	}
	return &domAcc
}

func NewTransactionService(repo domain.TransactionRepository, accountRepo domain.AccountRepository, manager db.TransactionManager) DefaultTransactionService {
	return DefaultTransactionService{Repository: repo, AccountRepository: accountRepo, TransactionManager: manager}
}
