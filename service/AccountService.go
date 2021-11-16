package service

import (
	"fmt"
	"github.com/barnettt/banking-lib/exceptions"
	"github.com/barnettt/banking-lib/logger"
	"github.com/barnettt/banking/domain"
	"github.com/barnettt/banking/dto"
	"github.com/golang-sql/civil"
	"time"
)

//go:generate mockgen -destination=../mock/service/mockAccountService.go -package=service github.com/barnettt/banking/service AccountService
type AccountService interface {
	Save(account *dto.AccountRequest) (*dto.NewAccountResponse, *exceptions.AppError)
	GetAccount(accountId string) (*domain.Account, *exceptions.AppError)
	UpdateAccount(account *domain.Account) *exceptions.AppError
}

type DefaultAccountService struct {
	repository domain.AccountRepository
}

func (defaultAccountService DefaultAccountService) Save(account *dto.AccountRequest) (*dto.NewAccountResponse, *exceptions.AppError) {
	dtoAcc, err := GetDomainAccount(account)
	if err != nil {
		logger.Error(err.Message)
		return nil, err
	}
	accResponse, err := defaultAccountService.repository.Save(*dtoAcc)
	if err != nil {
		logger.Error(err.Message)
		return nil, err
	}
	return accResponse, nil
}

func GetDomainAccount(request *dto.AccountRequest) (*domain.Account, *exceptions.AppError) {
	amount := fmt.Sprintf("%f", request.Amount)
	err := request.Validate()
	if err != nil {
		logger.Error("account validation failed : " + err.Message)
		return nil, err
	}
	// using time.Now().Format("2021-10-02T11:26:20")
	// kept failing as the date produced seemed to be mangled
	// using civil datetime seems to be ok with mysql db
	dateTime := civil.DateTimeOf(time.Now()).String()
	if request.OpeningDate != "" {
		dateTime = request.OpeningDate
	}
	domAcc := domain.Account{CustomerId: request.CustomerId,
		Amount:      amount,
		OpeningDate: dateTime,
		AccountType: request.AccountType,
		Status:      "1",
	}
	return &domAcc, nil
}

func (defaultAccountService DefaultAccountService) GetAccount(accountId string) (*domain.Account, *exceptions.AppError) {

	account, err := defaultAccountService.repository.GetAccount(accountId)
	if err != nil {
		return nil, err
	}
	domainAccount, appErr := GetDomainAccount(account)
	if err != nil {
		return nil, appErr
	}
	return domainAccount, nil
}

func (defaultAccountService DefaultAccountService) UpdateAccount(account *domain.Account) *exceptions.AppError {
	err := defaultAccountService.repository.UpdateAccount(*account)
	if err != nil {
		return err
	}
	return nil

}

func NewAccountService(repo domain.AccountRepository) DefaultAccountService {
	return DefaultAccountService{repository: repo}
}
